package fetch_items

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	item_info "github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	sq "github.com/Masterminds/squirrel"
)

type handler struct {
	itemRepo item_info.ReadRepository
}

func NewHandler(itemRepo item_info.ReadRepository) *handler {
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
			sq.Eq{"status": item_info.ActiveStatus},
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
		if r.URL.Query().Get("steam_gift") == "true" {
			criteries = append(criteries, sq.Eq{"is_steam_gift": true})
		} else {
			criteries = append(criteries, sq.Eq{"is_steam_gift": false})
		}
		if v := r.URL.Query().Get("region"); v != "" {
			decodedV, err := url.QueryUnescape(v)
			if err != nil {
				api.Return400("Параметр region должен быть в формате url encoded", w)
				return
			}
			log.Printf("trying to fetch items with region - %s", decodedV)
			criteries = append(criteries, sq.Eq{"region": decodedV})
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
			if _, ok := item_info.NotShowedItems[item.ID]; ok {
				continue
			}

			var oldPrice *float64
			if item.OldPrice != nil {
				oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
			}

			itemType := "cdkey"
			if item.IsSteamGift {
				itemType = "gift"
			}

			res = append(res, ItemDTO{
				ID:           item.ID,
				Title:        item.Title,
				Platform:     item.Platform,
				Region:       item.Region,
				CategoryID:   item.CategoryID,
				CurrentPrice: float64(item.CurrentPrice) / 100,
				IsForSale:    item.IsForSale,
				OldPrice:     oldPrice,
				ThumbnailURL: item.ThumbnailURL,
				Type:         itemType,
			})
		}
		api.ReturnOK(res, w)
	}
}
