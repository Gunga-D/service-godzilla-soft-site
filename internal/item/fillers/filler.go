package fillers

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type Filler interface {
	Fill(ctx context.Context, items []item.ItemCache) error
}
