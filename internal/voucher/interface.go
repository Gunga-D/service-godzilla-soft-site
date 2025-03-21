package voucher

import "context"

type Repository interface {
	CreateVoucher(ctx context.Context, v Voucher) (int64, error)
}
