package item

import "context"

type Repository interface {
	CreateItem(ctx context.Context, i Item) (int64, error)
}
