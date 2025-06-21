package create_order

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error)
}

type orderCreator interface {
	CreateItemOrder(ctx context.Context, email string, amount int64, itemID int64, itemSlip string, itemName string, utm *string) (string, error)
	CreateItemGiftOrder(ctx context.Context, steamProfile string, amount int64, itemID int64, utm *string) (string, error)
}

type userDatabus interface {
	PublishDatabusNewUserEmail(ctx context.Context, msg databus.NewUserEmailDTO) error
	PublishDatabusNewUserSteamLink(ctx context.Context, msg databus.NewUserSteamLinkDTO) error
}

type voucherActivation interface {
	ApplyVoucher(ctx context.Context, voucherValue string, i item.Item) (int64, error)
}
