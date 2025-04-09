package fetch_collections

import (
	"net/http"
	"strconv"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/collection"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	sq "github.com/Masterminds/squirrel"
)

type handler struct {
	repo collection.ReadRepository
}

func NewHandler(repo collection.ReadRepository) *handler {
	return &handler{
		repo: repo,
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

		criteries := sq.And{}
		if v, err := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64); err == nil {
			criteries = append(criteries, sq.Eq{"category_id": v})
		}
		collections, err := h.repo.FetchCollectionsByFilter(r.Context(), criteries, limit, offset)
		if err != nil {
			api.Return500("Ошибка получения подборок", w)
			return
		}
		res := make([]CollectionDTO, 0, len(collections))
		for _, col := range collections {
			res = append(res, CollectionDTO{
				ID:              col.ID,
				CategoryID:      col.CategoryID,
				Name:            col.Name,
				Description:     col.Description,
				BackgroundImage: col.BackgroundImage,
			})
		}

		api.ReturnOK(res, w)
	}
}
