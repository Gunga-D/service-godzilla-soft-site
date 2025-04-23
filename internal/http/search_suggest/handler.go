package search_suggest

import (
	"fmt"
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	suggester suggester
}

func NewHandler(suggester suggester) *handler {
	return &handler{
		suggester: suggester,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SearchSuggestRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		logger.Get().Log(fmt.Sprintf("🔍 Поиск товара по следующему запросу: \"%s\"", body.Query))

		suggested, err := h.suggester.Suggest(r.Context(), body.Query)
		if err != nil {
			api.Return500("Непредвиденная ошибка", w)
			return
		}

		res := make([]SearchSuggestDTO, 0, len(suggested))
		for _, s := range suggested {
			s := s
			if s.Type == "item" && s.Item != nil {
				var itemOldPrice *float64
				if s.Item.OldPrice != nil {
					itemOldPrice = pointer.ToFloat64(float64(*s.Item.OldPrice) / 100)
				}

				itemType := "cdkey"
				if s.Item.IsSteamGift {
					itemType = "gift"
				}

				itemHorizontalImage := s.Item.HorizontalImage
				if itemHorizontalImage == nil && s.Item.SteamBlock != nil {
					itemHorizontalImage = pointer.ToString(s.Item.SteamBlock.HeaderImage)
				}
				var itemGenres []string
				if s.Item.SteamBlock != nil {
					itemGenres = s.Item.SteamBlock.Genres
				}
				var itemReleaseDate *string
				if s.Item.SteamBlock != nil {
					itemReleaseDate = pointer.ToString(s.Item.SteamBlock.ReleaseDate)
				}

				res = append(res, SearchSuggestDTO{
					SuggestType:         s.Type,
					ItemID:              &s.Item.ID,
					ItemTitle:           &s.Item.Title,
					ItemCategoryID:      &s.Item.CategoryID,
					ItemCurrentPrice:    pointer.ToFloat64(float64(s.Item.CurrentPrice) / 100),
					ItemIsForSale:       &s.Item.IsForSale,
					ItemOldPrice:        itemOldPrice,
					ItemThumbnailURL:    &s.Item.ThumbnailURL,
					Probability:         s.Probability,
					ItemType:            pointer.ToString(itemType),
					ItemHorizontalImage: itemHorizontalImage,
					ItemGenres:          itemGenres,
					ItemReleaseDate:     itemReleaseDate,
				})
			}
			if s.Type == "banner" {
				res = append(res, SearchSuggestDTO{
					SuggestType: s.Type,
					BannerTitle: &s.Banner.Title,
					BannerImage: &s.Banner.Image,
					BannerURL:   &s.Banner.URL,
					Probability: s.Probability,
				})
			}
		}
		api.ReturnOK(SearchSuggestResponsePayload{
			Items: res,
		}, w)
	}
}
