package neuro_task

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/neuro"
)

type neuroSearch interface {
	Search(ctx context.Context, id string, query string) neuro.TaskResult
}

type neuroCache interface {
	SetTaskResult(ctx context.Context, id string, taskResult neuro.TaskResult) error
}
