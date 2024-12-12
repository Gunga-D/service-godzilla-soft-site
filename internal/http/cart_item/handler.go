package cart_item

import (
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	codeRepo              codeRepo
	itemGetter            itemGetter
	itemOutOfStockDatabus itemOutOfStockDatabus
}

func NewHandler(codeRepo codeRepo, itemGetter itemGetter, itemOutOfStockDatabus itemOutOfStockDatabus) *handler {
	return &handler{
		codeRepo:              codeRepo,
		itemGetter:            itemGetter,
		itemOutOfStockDatabus: itemOutOfStockDatabus,
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
			h.itemOutOfStockDatabus.PublishDatabusItemOutOfStockDTO(r.Context(), databus.ItemOutOfStockDTO{
				ItemID: body.ItemID,
			})

			api.Return409("Данный товар уже закончился", w)
			return
		}

		api.ReturnOK(CartItemResponsePayload{
			Price:    float64(item.CurrentPrice) / 100,
			Currency: "RUB",
		}, w)
	}
}
