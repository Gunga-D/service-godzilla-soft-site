package steam_invoice

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	orderCreator  orderCreator
	tinkoffClient tinkoff.Client
}

func NewHandler(orderCreator orderCreator, tinkoffClient tinkoff.Client) *handler {
	return &handler{
		orderCreator:  orderCreator,
		tinkoffClient: tinkoffClient,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SteamInvoiceRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		if body.Amount < 30 {
			api.Return400("Минимальная сумма от 30 рублей", w)
			return
		}
		if body.Amount > 30000 {
			api.Return400("Максимальная сумма до 30 000 рублей", w)
			return
		}
		_, err := url.ParseRequestURI(body.SteamLogin)
		if err == nil {
			api.Return400("Логин - это не ссылка на профиль пользователя", w)
			return
		}
		if body.SteamLogin == "" {
			api.Return400("Логин не может быть пустым", w)
			return
		}

		// В ордере сохраняем сумму, на которую пользователь хочет оплатить
		orderID, err := h.orderCreator.CreateSteamOrder(r.Context(), body.SteamLogin, body.Amount*100)
		if err != nil {
			log.Printf("[error] create order: %v", err)

			if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
				api.Return404("Данный товар закончился", w)
				return
			}

			api.Return500("Неизвестная ошибка", w)
			return
		}

		// Оплату же создаем на сумму + наша комиссия
		price := steam.SteamCalcPrice(body.Amount)
		invoiceResp, err := h.tinkoffClient.CreateInvoice(r.Context(), orderID, price*100, fmt.Sprintf("Steam пополнение на %d рублей", body.Amount))
		if err != nil {
			log.Printf("[error] cannot create invoice: %v", err)

			api.Return500("Неизвестная ошибка", w)
			return
		}

		logger.Get().Log(fmt.Sprintf("⚡️ На пополнение Steam на %d рублей создали ссылку на оплату", price))

		api.ReturnOK(SteamInvoiceResponse{
			OrderID:     orderID,
			PaymentLink: invoiceResp.PaymentURL,
		}, w)
	}
}
