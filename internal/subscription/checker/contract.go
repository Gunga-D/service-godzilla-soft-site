package checker

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription"
)

type subRepo interface {
	CreateSubscriptionBill(ctx context.Context, userID int64, amount int64, forPeriod time.Duration) (string, error)
	PaidSubscriptionBill(ctx context.Context, id string, rebillID string) error
	FetchLastUserSubscriptionBills(ctx context.Context) ([]subscription.PaidSubscription, error)
}
