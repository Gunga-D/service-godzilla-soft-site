package collection_details

import (
	"net/http"
	"strconv"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	getter getter
}

func NewHandler(getter getter) *handler {
	return &handler{
		getter: getter,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Query().Get("collection_id"), 10, 64)
		if err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		coll, err := h.getter.GetCollectionByID(r.Context(), id)
		if err != nil {
			api.Return500("Ошибка получения подборки", w)
			return
		}

		api.ReturnOK(CollectionDTO{
			ID:              coll.ID,
			CategoryID:      coll.CategoryID,
			Name:            coll.Name,
			Description:     coll.Description,
			BackgroundImage: coll.BackgroundImage,
		}, w)
	}
}
