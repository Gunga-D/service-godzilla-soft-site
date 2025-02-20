package new_items

import (
	"net/http"

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
		res := make([]ItemDTO, 0, len(newItems))
		for _, iID := range newItems {
			item, err := h.itemGetter.GetItemByID(r.Context(), iID)
			if err != nil {
				api.Return500("Неизвестная ошибка", w)
				return
			}
			if item == nil {
				continue
			}
			var oldPrice *float64
			if item.OldPrice != nil {
				oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
			}

			res = append(res, ItemDTO{
				ID:           item.ID,
				Title:        item.Title,
				Platform:     item.Platform,
				CategoryID:   item.CategoryID,
				Region:       item.Region,
				CurrentPrice: float64(item.CurrentPrice) / 100,
				IsForSale:    item.IsForSale,
				OldPrice:     oldPrice,
				ThumbnailURL: item.ThumbnailURL,
			})
		}
		api.ReturnOK(res, w)
	}
}
