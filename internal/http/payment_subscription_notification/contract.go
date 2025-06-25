package payment_subscription_notification

import "context"

type subRepo interface {
	PaidSubscriptionBill(ctx context.Context, id string, rebillID string) error
	FailedSubscriptionBill(ctx context.Context, id string) error
}
