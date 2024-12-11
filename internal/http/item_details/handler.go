package item_details

import (
	"net/http"
	"strconv"

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
		api.ReturnOK(ItemDTO{
			ID:           item.ID,
			Title:        item.Title,
			Description:  item.Description,
			CategoryID:   item.CategoryID,
			Platform:     item.Platform,
			Region:       item.Region,
			CurrentPrice: item.CurrentPrice,
			IsForSale:    item.IsForSale,
			OldPrice:     item.OldPrice,
			ThumbnailURL: item.ThumbnailURL,
			Slip:         item.Slip,
		}, w)
	}
}
