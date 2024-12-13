package delivery

import (
	"bytes"
	"context"
	"log"
	"text/template"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
)

type service struct {
	orderRepo        orderRepo
	yandexMailClient yandex_mail.Client
	deliveryTemplate *template.Template
}

func NewService(orderRepo orderRepo, yandexMailClient yandex_mail.Client, deliveryTemplate *template.Template) *service {
	return &service{
		orderRepo:        orderRepo,
		yandexMailClient: yandexMailClient,
		deliveryTemplate: deliveryTemplate,
	}
}

func (s *service) StartOrderDelivery(ctx context.Context) error {
	if err := s.process(ctx); err != nil {
		return err
	}

	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("order delivery stop")
			return nil
		case <-t.C:
			if err := s.process(ctx); err != nil {
				log.Printf("order delivery error: %v\n", err)
			}
		}
	}
}

func (s *service) process(ctx context.Context) error {
	orders, err := s.orderRepo.FetchPaidOrders(ctx)
	if err != nil {
		log.Printf("cannot process order delivery: fetch paid orders: %v\n", err)
		return err
	}

	for _, order := range orders {
		var body bytes.Buffer
		err = s.deliveryTemplate.Execute(&body, order)
		if err != nil {
			return err
		}

		err = s.yandexMailClient.SendMail([]string{
			order.Email,
		}, "Доставка товара от Godzilla Soft", body.String())
		if err != nil {
			return err
		}

		err = s.orderRepo.FinishOrder(ctx, order.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
