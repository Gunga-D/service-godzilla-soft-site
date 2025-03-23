package search_suggest

import (
	"fmt"
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	suggester  suggester
	itemGetter itemGetter
}

func NewHandler(suggester suggester, itemGetter itemGetter) *handler {
	return &handler{
		suggester:  suggester,
		itemGetter: itemGetter,
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

		items, err := h.suggester.Suggest(body.Query)
		if err != nil {
			api.Return500("Непредвиденная ошибка", w)
			return
		}

		res := make([]SearchSuggestDTO, 0, len(items))
		for _, item := range items {
			v, err := h.itemGetter.GetItemByName(r.Context(), item.Name)
			if err != nil {
				continue
			}
			if v == nil {
				continue
			}

			var itemOldPrice *float64
			if v.OldPrice != nil {
				itemOldPrice = pointer.ToFloat64(float64(*v.OldPrice) / 100)
			}

			res = append(res, SearchSuggestDTO{
				ItemID:           v.ID,
				ItemTitle:        v.Title,
				ItemCategoryID:   v.CategoryID,
				ItemCurrentPrice: float64(v.CurrentPrice) / 100,
				ItemIsForSale:    v.IsForSale,
				ItemOldPrice:     itemOldPrice,
				ItemThumbnailURL: v.ThumbnailURL,
				Probability:      item.Probability,
				ItemIsSteamGift:  v.IsSteamGift,
			})
		}
		api.ReturnOK(SearchSuggestResponsePayload{
			Items: res,
		}, w)
	}
}
