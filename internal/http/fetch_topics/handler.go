package fetch_topics

import (
	"fmt"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/cached"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

type Handler struct {
	repo *cached.Repo
}

type Response struct {
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

		fetchedTopics, err := h.repo.FetchTopics(r.Context(), limit, offset)
		if err != nil {
			api.Return500(fmt.Sprintf("Ошибка получения статей: %v", err), w)
			return
		}

		api.ReturnOK(toResponse(fetchedTopics), w)
	}
}

type topicResponse struct {
	Id         int64     `json:"id"`
	PreviewURL string    `json:"preview_url"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Convert your slice to []TopicResponse
func toResponse(topics []topics.Topic) []topicResponse {
	response := make([]topicResponse, len(topics))
	for i, t := range topics {
		response[i] = topicResponse{
			Id:         t.Id,
			PreviewURL: t.PreviewURL,
			Title:      t.Title,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
	}
	return response
}
