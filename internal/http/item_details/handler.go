package item_details

import (
	"net/http"
	"strconv"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	itemGetter itemGetter
}

func NewHandler(itemGetter itemGetter) *handler {
	return &handler{
		itemGetter: itemGetter,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Query().Get("item_id"), 10, 64)
		if err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		item, err := h.itemGetter.GetItemByID(r.Context(), id)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if item == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		var oldPrice *float64
		if item.OldPrice != nil {
			oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
		}

		api.ReturnOK(ItemDTO{
			ID:           item.ID,
			Title:        item.Title,
			Description:  item.Description,
			CategoryID:   item.CategoryID,
			Platform:     item.Platform,
			Region:       item.Region,
			CurrentPrice: float64(item.CurrentPrice) / 100,
			IsForSale:    item.IsForSale,
			OldPrice:     oldPrice,
			ThumbnailURL: item.ThumbnailURL,
			Slip:         item.Slip,
		}, w)
	}
}
