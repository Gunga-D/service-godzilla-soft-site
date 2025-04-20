package add_review

import (
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
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
		var body AddReviewRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		if body.Score > 5 || body.Score <= 0 {
			api.Return400("Невалидный запрос", w)
			return
		}

		if body.Comment != nil && len(*body.Comment) > 3000 {
			api.Return400("Невалидный запрос", w)
			return
		}

		var userID *int64
		if v, ok := r.Context().Value(user.MetaUserIDKey{}).(int64); ok {
			userID = pointer.ToInt64(v)
		}

		reviewID, err := h.repo.AddReview(r.Context(), userID, body.ItemID, body.Comment, body.Score)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(AddReviewResponse{
			ReviewID: reviewID,
		}, w)
	}
}
