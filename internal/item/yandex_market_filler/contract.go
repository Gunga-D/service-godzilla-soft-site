package yandex_market_filler

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemGetter interface {
	FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error)
}
