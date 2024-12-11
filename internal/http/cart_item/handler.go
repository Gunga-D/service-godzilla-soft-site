package cart_item

import (
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	codeRepo   codeRepo
	itemGetter itemGetter
}

func NewHandler(codeRepo codeRepo, itemGetter itemGetter) *handler {
	return &handler{
		codeRepo:   codeRepo,
		itemGetter: itemGetter,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CartItemRequest
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

		hasCodes, err := h.codeRepo.HasActiveCode(r.Context(), body.ItemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if !hasCodes {
			api.Return409("Данный товар уже закончился", w)
			return
		}

		api.ReturnOK(CartItemResponsePayload{
			Price:    item.CurrentPrice,
			Currency: "RUB",
			// TODO: Добавить paymentLink
			PaymentLink: "",
		}, w)
	}
}
