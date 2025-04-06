package search_suggest

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/suggest"
)

type suggester interface {
	Suggest(ctx context.Context, text string) ([]suggest.Suggested, error)
}
