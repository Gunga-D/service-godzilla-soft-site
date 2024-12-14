package search_suggest

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/suggest"
)

type suggester interface {
	Suggest(text string) ([]suggest.SuggestedItem, error)
}

type itemGetter interface {
	GetItemByName(ctx context.Context, name string) (*item.Item, error)
}
