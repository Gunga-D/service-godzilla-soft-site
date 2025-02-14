package cart_item

import (
	"fmt"
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
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

		i, err := h.itemGetter.GetItemByID(r.Context(), body.ItemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if i == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		hasCodes, err := h.codeRepo.HasActiveCode(r.Context(), i.ID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if !hasCodes {
			h.itemOutOfStockDatabus.PublishDatabusChangeItemState(r.Context(), databus.ChangeItemStateDTO{
				ItemID: body.ItemID,
				Status: item.PausedStatus,
			})

			logger.Get().Log(fmt.Sprintf("❗️ Товар\"%s\" закончился", i.Title))

			api.Return409("Данный товар уже закончился", w)
			return
		}

		api.ReturnOK(CartItemResponsePayload{
			Price:    float64(i.CurrentPrice) / 100,
			Currency: "RUB",
		}, w)
	}
}
