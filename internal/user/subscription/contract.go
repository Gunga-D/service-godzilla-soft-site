package subscription

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription"
)

type subRepo interface {
	GetLastUserSubscriptionBill(ctx context.Context, userID int64) (*subscription.UserSubscription, error)
}
