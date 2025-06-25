package user_profile

import "context"

type subDeterm interface {
	HasSubscription(ctx context.Context, userID int64) (bool, error)
}
