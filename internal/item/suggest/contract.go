package suggest

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type getter interface {
	FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error)
}

type itemCache interface {
	GetItemByName(ctx context.Context, name string) (*item.ItemCache, error)
}
