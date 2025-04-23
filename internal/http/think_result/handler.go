package think_result

import (
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	thinker thinker
}

func NewHandler(thinker thinker) *handler {
	return &handler{
		thinker: thinker,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ThinkResultRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		res, err := h.thinker.GetThinkingResult(r.Context(), body.ID)
		if err != nil {
			api.Return500("Непредвиденная ошибка", w)
			return
		}
		if res == nil {
			api.Return404("Еще не обработана", w)
			return
		}
		items := make([]ItemDTO, 0, len(res.Items))
		for _, i := range res.Items {
			itemType := "cdkey"
			if i.IsSteamGift {
				itemType = "gift"
			}

			itemHorizontalImage := i.HorizontalImage
			if itemHorizontalImage == nil && i.SteamBlock != nil {
				itemHorizontalImage = pointer.ToString(i.SteamBlock.HeaderImage)
			}
			var itemGenres []string
			if i.SteamBlock != nil {
				itemGenres = i.SteamBlock.Genres
			}
			var itemReleaseDate *string
			if i.SteamBlock != nil {
				itemReleaseDate = pointer.ToString(i.SteamBlock.ReleaseDate)
			}

			items = append(items, ItemDTO{
				ID:              i.ID,
				CategoryID:      i.CategoryID,
				Title:           i.Title,
				CurrentPrice:    float64(i.CurrentPrice) / 100,
				ThumbnailURL:    i.ThumbnailURL,
				Type:            itemType,
				HorizontalImage: itemHorizontalImage,
				Genres:          itemGenres,
				ReleaseDate:     itemReleaseDate,
			})
		}
		api.ReturnOK(ThinkResultResponse{
			Reflection: res.Reflection,
			Items:      items,
		}, w)
	}
}
