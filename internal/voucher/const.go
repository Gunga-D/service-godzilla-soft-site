package voucher

import "errors"

const (
	DirectAmountType   = "direct"
	FloatingAmountType = "floating"
)

var (
	ErrNotFoundVoucher = errors.New("voucher not found")
)
