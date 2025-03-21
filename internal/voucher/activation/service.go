package activation

import (
	"context"
	"errors"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/voucher"
)

type service struct {
	voucherRepo voucherRepo
}

func NewService(voucherRepo voucherRepo) *service {
	return &service{
		voucherRepo: voucherRepo,
	}
}

func (s *service) ApplyVoucher(ctx context.Context, voucherValue string, i item.Item) (int64, error) {
	v, err := s.voucherRepo.ApplyVoucher(ctx, voucherValue)
	if err != nil {
		return 0, err
	}
	newPrice, _, err := s.calcVoucher(*v, i)
	return newPrice, err
}

func (s *service) PeekVoucher(ctx context.Context, voucherValue string, i item.Item) (int64, bool, error) {
	v, err := s.voucherRepo.GetActiveVoucherByValue(ctx, voucherValue)
	if err != nil {
		return 0, false, err
	}
	if v == nil {
		return 0, false, voucher.ErrNotFoundVoucher
	}
	return s.calcVoucher(*v, i)
}

func (s *service) calcVoucher(v voucher.Voucher, i item.Item) (int64, bool, error) {
	newAmount := int64(0)
	switch v.Type {
	case voucher.DirectAmountType:
		newAmount = i.CurrentPrice - v.Impact
	case voucher.FloatingAmountType:
		newAmount = i.CurrentPrice - (i.CurrentPrice*v.Impact)/100
	default:
		return 0, false, errors.New("type of voucher is not supported")
	}
	if i.LimitPrice != nil {
		if newAmount < *i.LimitPrice {
			return *i.LimitPrice, true, nil
		}
	}
	return newAmount, false, nil
}
