package fetch_items

import (
	"net/http"
	"strconv"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	sq "github.com/Masterminds/squirrel"
)

type handler struct {
	itemRepo item.ReadRepository
}

func NewHandler(itemRepo item.ReadRepository) *handler {
	return &handler{
		itemRepo: itemRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = 10
		}
		offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			offset = 0
		}

		criteries := sq.And{
			sq.Eq{"status": item.ActiveStatus},
		}
		if v, err := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64); err == nil {
			criteries = append(criteries, sq.GtOrEq{"current_price": int64(v * 100)})
		}
		if v, err := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64); err == nil {
			criteries = append(criteries, sq.LtOrEq{"current_price": int64(v * 100)})
		}
		if v, err := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64); err == nil {
			criteries = append(criteries, sq.Eq{"category_id": v})
		}
		if v := r.URL.Query().Get("region"); v != "" {
			criteries = append(criteries, sq.Eq{"region": v})
		}
		if v := r.URL.Query().Get("platform"); v != "" {
			criteries = append(criteries, sq.Eq{"platform": v})
		}

		items, err := h.itemRepo.FetchItemsByFilter(r.Context(), criteries, limit, offset)
		if err != nil {
			api.Return500("Ошибка получения каталога", w)
			return
		}
		res := make([]ItemDTO, 0, len(items))
		for _, item := range items {
			var oldPrice *float64
			if item.OldPrice != nil {
				oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
			}

			res = append(res, ItemDTO{
				ID:           item.ID,
				Title:        item.Title,
				Platform:     item.Platform,
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
