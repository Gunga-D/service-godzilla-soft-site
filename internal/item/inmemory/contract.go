package inmemory

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type getter interface {
	FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error)
}

type recomendation interface {
	Sync(ctx context.Context, itemsBySteamAppID map[int64]item.ItemCache) error
}
