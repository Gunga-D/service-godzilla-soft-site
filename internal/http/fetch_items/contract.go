package fetch_items

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemCache interface {
	GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error)
}
