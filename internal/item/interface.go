package item

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type WriteRepository interface {
	CreateItem(ctx context.Context, i Item) (int64, error)
}

type ReadRepository interface {
	FetchItemsByFilter(ctx context.Context, criteria sq.And, limit uint64, offset uint64) ([]Item, error)
}
