package subscription

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription"
)

type service struct {
	subRepo subRepo
}

func NewService(subRepo subRepo) *service {
	return &service{
		subRepo: subRepo,
	}
}

func (s *service) HasSubscription(ctx context.Context, userID int64) (bool, error) {
	sub, err := s.subRepo.GetLastUserSubscriptionBill(ctx, userID)
	if err != nil {
		return false, err
	}
	if sub == nil {
		return false, nil
	}
	if sub.Status != subscription.PaidStatus {
		return false, nil
	}
	if time.Now().Unix() > sub.ExpiredAt {
		return false, nil
	}
	return true, nil
}
