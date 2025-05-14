package think_result

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/neuro"
)

type neuroCache interface {
	GetTaskResult(ctx context.Context, id string) (*neuro.TaskResult, error)
}
