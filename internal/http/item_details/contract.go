package item_details

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error)
}

type recomendation interface {
	Recommend(ctx context.Context, itemID int64, genres []string) ([]item.ItemCache, error)
}
