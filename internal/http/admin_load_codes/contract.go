package admin_load_codes

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemRepo interface {
	GetItemByID(ctx context.Context, id int64) (*item.Item, error)
}

type codeRepo interface {
	CreateCodes(ctx context.Context, itemID int64, value []string) error
}

type itemChangeStateDatabus interface {
	PublishDatabusChangeItemState(ctx context.Context, msg databus.ChangeItemStateDTO) error
}
