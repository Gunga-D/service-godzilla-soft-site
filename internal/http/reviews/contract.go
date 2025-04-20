package reviews

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/review"
)

type repo interface {
	GetScore(ctx context.Context, itemID int64) (float64, error)
	FetchCommentReviews(ctx context.Context, itemID int64, limit uint64, offset uint64) ([]review.Review, error)
}
