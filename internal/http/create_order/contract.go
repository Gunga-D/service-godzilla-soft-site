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
	CreateItemOrder(ctx context.Context, email string, amount int64, itemID int64, itemSlip string) (string, error)
}

type userRegistrationDatabus interface {
	PublishDatabusQuickUserRegistration(ctx context.Context, msg databus.QuickUserRegistrationDTO) error
}

type voucherActivation interface {
	ApplyVoucher(ctx context.Context, voucherValue string, i item.Item) (int64, error)
}
