package checker

import (
	"context"
	"log"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
)

const (
	_subscriptionCost     = 25000
	_subscriptionDuration = time.Hour * 24 * 31
)

type service struct {
	subRepo       subRepo
	tinkoffClient tinkoff.Client
}

func NewService(subRepo subRepo, tinkoffClient tinkoff.Client) *service {
	return &service{
		subRepo:       subRepo,
		tinkoffClient: tinkoffClient,
	}
}

func (s *service) StartCheck(ctx context.Context) error {
	log.Println("start up checking subscription")
	if err := s.process(ctx); err != nil {
		log.Printf("subscription check error: %v", err)
		return err
	}

	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("subscription check stop")
			return nil
		case <-t.C:
			if err := s.process(ctx); err != nil {
				log.Printf("subscription check error: %v\n", err)
			}
		}
	}
}

func (s *service) process(ctx context.Context) error {
	subs, err := s.subRepo.FetchLastUserSubscriptionBills(ctx)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if !sub.NeedProlong {
			continue
		}
		if time.Now().Unix() < sub.ExpiredAt {
			continue
		}

		id, err := s.subRepo.CreateSubscriptionBill(ctx, sub.UserID, _subscriptionCost, time.Duration(sub.ExpiredAt-sub.CreatedAt)*time.Second)
		if err != nil {
			log.Printf("[error] cannot create subscription bill: %v\n", err)
			continue
		}
		paymentResp, err := s.tinkoffClient.CreateInvoice(ctx, id, _subscriptionCost, "Подписка GODZILLA SOFT на библиотеку Steam игр")
		if err != nil {
			log.Printf("[error] cannot create subscription invoice: %v\n", err)
			continue
		}
		if !paymentResp.Success {
			log.Printf("[error] cannot create subscription invoice: [error code - %v] [status - %v]\n", paymentResp.ErrorCode, paymentResp.Status)
			continue
		}

		chargeResp, err := s.tinkoffClient.Charge(ctx, sub.RebillID, paymentResp.PaymentId)
		if err != nil {
			log.Printf("[error] cannot create subscription charge: %v\n", err)
			continue
		}
		if !chargeResp.Success {
			log.Printf("[error] cannot create subscription charge: [error code - %v] [status - %v]\n", chargeResp.ErrorCode, chargeResp.Status)
			continue
		}
		err = s.subRepo.PaidSubscriptionBill(ctx, id, sub.RebillID)
		if err != nil {
			log.Printf("[error] cannot subscription change the status to paid: %v\n", err)
			continue
		}
	}
	return nil
}
