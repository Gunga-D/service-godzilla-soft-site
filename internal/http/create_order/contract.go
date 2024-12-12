package create_order

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.Item, error)
}

type orderCreator interface {
	CreateOrder(ctx context.Context, email string, amount int64, itemID int64) (string, error)
}

type userRegistrationDatabus interface {
	PublishDatabusQuickUserRegistration(ctx context.Context, msg databus.QuickUserRegistrationDTO) error
}
