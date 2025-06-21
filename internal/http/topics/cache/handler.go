package cache

import (
	"fmt"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/cached"
	"net/http"
)

type Handler struct {
	repo *cached.Repo
}

func NewHandler(repo *cached.Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.repo.InvalidateCache(r.Context())
		if err != nil {
			api.Return500(fmt.Sprintf("Ошибка очистки кэша статей: %v", err), w)
			return
		}

		api.ReturnOK(nil, w)
	}
}

func (h *Handler) HandleSync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.repo.Sync(r.Context())
		if err != nil {
			api.Return500(fmt.Sprintf("Ошибка синхронизации кэша статей: %v", err), w)
			return
		}

		api.ReturnOK(nil, w)
	}
}
