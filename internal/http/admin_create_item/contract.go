package admin_create_item

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemsRepo interface {
	CreateItem(ctx context.Context, i item.Item) (int64, error)
}
