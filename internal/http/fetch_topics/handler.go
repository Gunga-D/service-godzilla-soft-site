package fetch_topics

import (
	"fmt"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/cached"
	"net/http"
	"strconv"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

type Handler struct {
	repo *cached.Repo
}

func NewHandler(repo *cached.Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = defaultLimit
		}
		offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			offset = defaultOffset
		}

		previews, err := h.repo.FetchTopics(r.Context(), limit, offset)
		if err != nil {
			api.Return500(fmt.Sprintf("Ошибка получения статей: %v", err), w)
			return
		}

		api.ReturnOK(previews, w)
	}
}
