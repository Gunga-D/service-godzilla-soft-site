package subscription_product

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription"
)

type userChecker interface {
	HasSubscription(ctx context.Context, userID int64) (bool, error)
}

type subRepo interface {
	GetSubscriptionProduct(ctx context.Context, itemID int64) (*subscription.SubscriptionProduct, error)
}
