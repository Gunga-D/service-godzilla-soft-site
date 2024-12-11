package cart_item

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type codeRepo interface {
	HasActiveCode(ctx context.Context, itemID int64) (bool, error)
}

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.Item, error)
}
