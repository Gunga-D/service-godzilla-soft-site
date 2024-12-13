package delivery

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/order"
)

type orderRepo interface {
	FetchPaidOrders(ctx context.Context) ([]order.PaidOrder, error)
	FinishOrder(ctx context.Context, orderID string) error
}
