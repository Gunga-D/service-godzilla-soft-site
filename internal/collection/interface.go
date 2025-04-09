package collection

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type ReadRepository interface {
	FetchCollectionsByFilter(ctx context.Context, criteria sq.And, limit uint64, offset uint64) ([]Collection, error)
	FetchCollectionItemsByCollectionID(ctx context.Context, collectionID int64, limit uint64, offset uint64) ([]CollectionItem, error)
}
