package fetch_items

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	item_info "github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	sq "github.com/Masterminds/squirrel"
)

const (
	_defaultCount = 0
)

type Handler struct {
	itemRepo  item_info.ReadRepository
	itemCache itemCache
}

func NewHandler(itemRepo item_info.ReadRepository, itemCache itemCache) *Handler {
	return &Handler{
		itemRepo:  itemRepo,
		itemCache: itemCache,
	}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = 10
		}
		offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			offset = 0
		}

		orderBy := []string{
			"id",
		}
		if r.URL.Query().Get("random") == "1" {
			orderBy = []string{
				"random()",
			}
		}
		if r.URL.Query().Get("popular") == "1" {
			orderBy = []string{
				"popular",
				"id",
			}
		}
		if r.URL.Query().Get("new") == "1" {
			orderBy = []string{
				"new",
				"id",
			}
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
		if r.URL.Query().Get("in_sub") == "1" {
			criteries = append(criteries, sq.Eq{"in_sub": true})
		} else {
			if r.URL.Query().Get("steam_gift") == "true" {
				criteries = append(criteries, sq.Eq{"is_steam_gift": true})
			} else {
				criteries = append(criteries, sq.Eq{"is_steam_gift": false})
			}
		}
		if v := r.URL.Query().Get("region"); v != "" {
			regionCriteria := sq.Or{}
			for _, rawRegion := range strings.Split(v, ";") {
				decodedRegion, err := url.QueryUnescape(rawRegion)
				if err != nil {
					api.Return400("Параметр region должен быть в формате url encoded", w)
					return
				}
				regionCriteria = append(regionCriteria, sq.Eq{"region": decodedRegion})
			}
			criteries = append(criteries, regionCriteria)
		}
		if v := r.URL.Query().Get("platform"); v != "" {
			platformCriteria := sq.Or{}
			for _, platform := range strings.Split(v, ";") {
				platformCriteria = append(platformCriteria, sq.Eq{"platform": platform})
			}
			criteries = append(criteries, platformCriteria)
		}
		if r.URL.Query().Get("unavailable") == "1" {
			criteries = append(criteries, sq.Eq{"unavailable": true})
		}
		items, err := h.itemRepo.FetchItemsByFilter(r.Context(), criteries, limit, offset, orderBy)
		if err != nil {
			api.Return500("Ошибка получения каталога", w)
			return
		}
		itemsCount, err := h.itemRepo.GetItemsCountByFilter(r.Context(), criteries)
		if err != nil {
			log.Println(err)
			itemsCount = _defaultCount
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

			cacheItem, err := h.itemCache.GetItemByID(r.Context(), item.ID)
			if err != nil {
				continue
			}
			if cacheItem == nil {
				continue
			}

			desc := item.Description
			if desc == nil && cacheItem.SteamBlock != nil {
				desc = pointer.ToString(cacheItem.SteamBlock.ShortDescription)
			}

			var genres []string
			if cacheItem.SteamBlock != nil {
				genres = cacheItem.SteamBlock.Genres
			}

			horizontalImageURL := item.HorizontalImage
			if horizontalImageURL == nil && cacheItem.SteamBlock != nil {
				horizontalImageURL = pointer.ToString(cacheItem.SteamBlock.HeaderImage)
			}

			var releaseDate *string
			if cacheItem.SteamBlock != nil {
				releaseDate = pointer.ToString(cacheItem.SteamBlock.ReleaseDate)
			}

			res = append(res, ItemDTO{
				ID:                 item.ID,
				Title:              item.Title,
				Platform:           item.Platform,
				Region:             item.Region,
				CategoryID:         item.CategoryID,
				CurrentPrice:       float64(item.CurrentPrice) / 100,
				IsForSale:          item.IsForSale,
				OldPrice:           oldPrice,
				ThumbnailURL:       item.ThumbnailURL,
				Type:               itemType,
				Description:        desc,
				Genres:             genres,
				HorizontalImageURL: horizontalImageURL,
				ReleaseDate:        releaseDate,
				TotalCount:         &itemsCount,
				InSub:              item.InSub,
			})
		}
		api.ReturnOK(res, w)
	}
}
