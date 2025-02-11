package payment_notification

import "context"

type orderRepo interface {
	PaidOrder(ctx context.Context, orderID string) error
	FailedOrder(ctx context.Context, orderID string) error
}
