package create_order

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/voucher"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	itemGetter              itemGetter
	orderCreator            orderCreator
	tinkoffClient           tinkoff.Client
	userRegistrationDatabus userRegistrationDatabus
	voucherActivation       voucherActivation
}

func NewHandler(itemGetter itemGetter,
	orderCreator orderCreator,
	tinkoffClient tinkoff.Client,
	userRegistrationDatabus userRegistrationDatabus,
	voucherActivation voucherActivation,
) *handler {

	return &handler{
		itemGetter:              itemGetter,
		orderCreator:            orderCreator,
		tinkoffClient:           tinkoffClient,
		userRegistrationDatabus: userRegistrationDatabus,
		voucherActivation:       voucherActivation,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateOrderRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		item, err := h.itemGetter.GetItemByID(r.Context(), body.ItemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if item == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		currentPrice := item.CurrentPrice
		if body.Voucher != nil {
			currentPrice, err = h.voucherActivation.ApplyVoucher(r.Context(), *body.Voucher, item.Item)
			if err != nil {
				if errors.Is(err, voucher.ErrNotFoundVoucher) {
					api.Return404("Купон не найден или уже был активирован", w)
					return
				}
				api.Return500("Непредвиденная ошибка", w)
				return
			}
		}

		var orderID string
		if item.IsSteamGift {
			// Отдельное флоу для Стим гифта
			if body.SteamProfile == nil {
				api.Return400("Для товара \"Steam gift\" поле Steam профиль обязательно для заполнения", w)
				return
			}

			giftOrderID, err := h.orderCreator.CreateItemGiftOrder(r.Context(), *body.SteamProfile, currentPrice, item.ID)
			if err != nil {
				log.Printf("[error] create gift order: %v", err)
				api.Return500("Неизвестная ошибка", w)
				return
			}
			orderID = giftOrderID
		} else {
			// Флоу для ключа
			var userEmail string
			if body.Email != nil {
				userEmail = *body.Email

				if ok := auth.ValidateEmail(userEmail); !ok {
					api.Return400("Некорректная почта, исправьте ее и попробуйте еще раз", w)
					return
				}

				err = h.userRegistrationDatabus.PublishDatabusQuickUserRegistration(r.Context(), databus.QuickUserRegistrationDTO{
					Email: userEmail,
				})
				if err != nil {
					api.Return500("Непредвиденная ошибка во время быстрой регистрации пользователя", w)
					return
				}
			} else {
				email, ok := r.Context().Value(user.MetaUserEmailKey{}).(string)
				if !ok {
					api.Return401("Ошибка авторизации", w)
					return
				}
				userEmail = email
			}

			cdkeyOrderID, err := h.orderCreator.CreateItemOrder(r.Context(), userEmail, currentPrice, item.ID, item.Slip)
			if err != nil {
				log.Printf("[error] create order: %v", err)

				if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
					api.Return404("Данный товар закончился", w)
					return
				}

				api.Return500("Неизвестная ошибка", w)
				return
			}
			orderID = cdkeyOrderID
		}

		invoiceResp, err := h.tinkoffClient.CreateInvoice(r.Context(), orderID, currentPrice, fmt.Sprintf("Покупка \"%s\"", item.Title))
		if err != nil {
			log.Printf("[error] cannot create invoice: %v", err)

			api.Return500("Неизвестная ошибка", w)
			return
		}

		logger.Get().Log(fmt.Sprintf("⚡️ На товар\"%s\" создали ссылку на оплату", item.Title))

		api.ReturnOK(CreateOrderResponsePayload{
			OrderID:     orderID,
			PaymentLink: invoiceResp.PaymentURL,
		}, w)
	}
}
