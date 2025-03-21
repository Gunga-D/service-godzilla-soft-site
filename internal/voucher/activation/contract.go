package activation

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/voucher"
)

type voucherRepo interface {
	ApplyVoucher(ctx context.Context, value string) (*voucher.Voucher, error)
	GetActiveVoucherByValue(ctx context.Context, value string) (*voucher.Voucher, error)
}
