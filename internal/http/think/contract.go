package think

import (
	"context"
)

type thinker interface {
	StartThinking(ctx context.Context, query string) string
}
