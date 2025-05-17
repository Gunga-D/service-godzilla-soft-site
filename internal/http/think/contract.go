package think

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type neuroDatabus interface {
	PublishDatabusNeuroTask(ctx context.Context, msg databus.NeuroTaskDTO) error
	PublishDatabusNeuroNewItems(ctx context.Context, msg databus.NeuroNewItemsDTO) error
}
