package think_result

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/deepthink"
)

type thinker interface {
	GetThinkingResult(ctx context.Context, id string) (*deepthink.ThinkResult, error)
}
