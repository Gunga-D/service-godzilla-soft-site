package check_voucher

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error)
}

type voucherActivation interface {
	PeekVoucher(ctx context.Context, voucherValue string, i item.Item) (int64, bool, error)
}
