package add_review

import "context"

type repo interface {
	AddReview(ctx context.Context, userID *int64, itemID int64, comment *string, score int) (int64, error)
}
