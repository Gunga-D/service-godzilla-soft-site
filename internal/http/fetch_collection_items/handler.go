package fetch_collection_items

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/collection"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	itemGetter itemGetter
	repo       collection.ReadRepository
}

func NewHandler(itemGetter itemGetter, repo collection.ReadRepository) *handler {
	return &handler{
		itemGetter: itemGetter,
		repo:       repo,
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

		var collectionID int64
		if v, err := strconv.ParseInt(r.URL.Query().Get("collection_id"), 10, 64); err == nil {
			collectionID = v
		} else {
			api.Return400("Обязательно необходимо указать id подборки", w)
			return
		}
		colItems, err := h.repo.FetchCollectionItemsByCollectionID(r.Context(), collectionID, limit, offset)
		if err != nil {
			log.Printf("[error] cannot fetch collection items: %v\n", err)
			api.Return500("Ошибка получения айтемов из подборки", w)
			return
		}
		res := make([]CollectionItemDTO, 0, len(colItems))
		for _, coli := range colItems {
			i, err := h.itemGetter.GetItemByID(r.Context(), coli.ItemID)
			if err != nil {
				api.Return500("Непредвиденная ошибка", w)
				return
			}

			var oldPrice *float64
			if i.OldPrice != nil {
				oldPrice = pointer.ToFloat64(float64(*i.OldPrice) / 100)
			}

			itemType := "cdkey"
			if i.IsSteamGift {
				itemType = "gift"
			}

			horizontalImage := i.HorizontalImage
			if horizontalImage == nil && i.SteamBlock != nil {
				horizontalImage = pointer.ToString(i.SteamBlock.HeaderImage)
			}
			var releaseDate *string
			if i.SteamBlock != nil {
				releaseDate = pointer.ToString(i.SteamBlock.ReleaseDate)
			}
			var genres []string
			if i.SteamBlock != nil {
				genres = i.SteamBlock.Genres
			}

			res = append(res, CollectionItemDTO{
				ID:                 i.ID,
				Title:              i.Title,
				CategoryID:         i.CategoryID,
				Platform:           i.Platform,
				Region:             i.Region,
				CurrentPrice:       float64(i.CurrentPrice) / 100,
				IsForSale:          i.IsForSale,
				OldPrice:           oldPrice,
				ThumbnailURL:       i.ThumbnailURL,
				HorizontalImageURL: horizontalImage,
				Type:               itemType,
				ReleaseDate:        releaseDate,
				Genres:             genres,
			})
		}

		api.ReturnOK(res, w)
	}
}
