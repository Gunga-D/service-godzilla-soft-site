package think_result

import (
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	neuroCache neuroCache
}

func NewHandler(neuroCache neuroCache) *handler {
	return &handler{
		neuroCache: neuroCache,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ThinkResultRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		taskRes, err := h.neuroCache.GetTaskResult(r.Context(), body.ID)
		if err != nil {
			api.Return500("Непредвиденная ошибка", w)
			return
		}
		if taskRes == nil {
			api.Return404("Еще не обработана", w)
			return
		}

		if !taskRes.Success || taskRes.Data == nil {
			if taskRes.Message != nil {
				api.ReturnOK(ThinkResultResponse{
					Reflection: *taskRes.Message,
					Items:      []ItemDTO{},
				}, w)
				return
			}
			api.Return500("Непредвиденная ошибка", w)
			return
		}

		items := make([]ItemDTO, 0, len(taskRes.Data.Items))
		for _, i := range taskRes.Data.Items {
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
			Reflection: taskRes.Data.Reflection,
			Items:      items,
		}, w)
	}
}
