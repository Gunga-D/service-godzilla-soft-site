package reviews

import (
	"fmt"
	"net/http"
	"strconv"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	repo repo
}

func NewHandler(repo repo) *handler {
	return &handler{
		repo: repo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := strconv.ParseInt(r.URL.Query().Get("item_id"), 10, 64)
		if err != nil {
			api.Return400("Обязательный параметр - item_id", w)
			return
		}
		limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = 10
		}
		offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			offset = 0
		}
		score, err := h.repo.GetScore(r.Context(), itemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		reviews, err := h.repo.FetchCommentReviews(r.Context(), itemID, limit, offset)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		resReviews := make([]ReviewDTO, 0, len(reviews))
		for _, review := range reviews {
			resReviews = append(resReviews, ReviewDTO{
				Comment:   review.Comment,
				Score:     review.Score,
				CreatedAt: fmt.Sprintf("%d %s %d", review.CreatedAt.Day(), RussianMonth(review.CreatedAt), review.CreatedAt.Year()),
			})
		}

		api.ReturnOK(ReviewsResponse{
			Score:   Round(score, 2),
			Reviews: resReviews,
		}, w)
	}
}
