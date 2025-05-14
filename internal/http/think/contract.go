package think

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type neuroTaskDatabus interface {
	PublishDatabusNeuroTask(ctx context.Context, msg databus.NeuroTaskDTO) error
}
