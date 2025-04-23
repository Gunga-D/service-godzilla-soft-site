package deepthink

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemCache interface {
	GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error)
}

type getter interface {
	FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error)
}
