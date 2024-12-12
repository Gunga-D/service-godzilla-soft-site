package create_order

import (
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type handler struct {
	itemGetter              itemGetter
	orderCreator            orderCreator
	userRegistrationDatabus userRegistrationDatabus
}

func NewHandler(itemGetter itemGetter, orderCreator orderCreator, userRegistrationDatabus userRegistrationDatabus) *handler {
	return &handler{
		itemGetter:              itemGetter,
		orderCreator:            orderCreator,
		userRegistrationDatabus: userRegistrationDatabus,
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

		var userEmail string
		if body.Email != nil {
			userEmail = *body.Email

			err = h.userRegistrationDatabus.PublishDatabusQuickUserRegistration(r.Context(), databus.QuickUserRegistrationDTO{
				Email: userEmail,
			})
			if err != nil {
				api.Return500("Непредвиденная ошибка во время быстрой регистрации пользователя", w)
				return
			}
		} else {
			userEmail = r.Context().Value(user.MetaUserEmailKey{}).(string)
		}

		orderID, err := h.orderCreator.CreateOrder(r.Context(), userEmail, item.CurrentPrice, item.ID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(CreateOrderResponsePayload{
			OrderID:     orderID,
			PaymentLink: "",
		}, w)
	}
}
