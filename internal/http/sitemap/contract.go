package sitemap

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type getter interface {
	FetchAllItems(ctx context.Context) ([]item.ItemCache, error)
}
