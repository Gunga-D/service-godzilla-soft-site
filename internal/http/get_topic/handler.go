package get_topic

import (
	"fmt"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/cached"
	"net/http"
	"strconv"
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
		id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			api.Return400("Couldn't parse id query parameter", w)
			return
		}

		topic, err := h.repo.GetTopic(r.Context(), id)
		if err != nil {
			api.Return500(fmt.Sprintf("Ошибка при получении статьи: %v", err), w)
			return
		}
		api.ReturnOK(topic, w)
	}
}
