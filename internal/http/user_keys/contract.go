package user_keys

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/order"
)

type orderRepo interface {
	FetchUserOrdersByEmail(ctx context.Context, email string) ([]order.UserOrder, error)
}
