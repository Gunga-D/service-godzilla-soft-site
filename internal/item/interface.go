package item

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type WriteRepository interface {
	CreateItem(ctx context.Context, i Item) (int64, error)
	UpdatePrice(ctx context.Context, itemID int64, currentPrice int64, limitPrice int64, priceLoc string) error
	UpdateSteamRawData(ctx context.Context, itemID int64, steamRawData string) error
}

type ReadRepository interface {
	GetItemsCountByFilter(ctx context.Context, criteria sq.And) (int64, error)
	FetchItemsByFilter(ctx context.Context, criteria sq.And, limit uint64, offset uint64, orderBy []string) ([]Item, error)
}

type Repository interface {
	WriteRepository
	ReadRepository
}
