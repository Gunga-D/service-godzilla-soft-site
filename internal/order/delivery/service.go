package delivery

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam_invoice"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
)

type service struct {
	orderRepo          orderRepo
	steamInvoiceClient steam_invoice.Client
	yandexMailClient   yandex_mail.Client
	deliveryTemplate   *template.Template
}

func NewService(orderRepo orderRepo, steamInvoiceClient steam_invoice.Client, yandexMailClient yandex_mail.Client, deliveryTemplate *template.Template) *service {
	return &service{
		orderRepo:          orderRepo,
		steamInvoiceClient: steamInvoiceClient,
		yandexMailClient:   yandexMailClient,
		deliveryTemplate:   deliveryTemplate,
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
		if order.CodeValue == "STEAM_INVOICE_BY_LOGIN" {
			err := s.steamInvoiceClient.CreateInvoice(ctx, order.Email, order.Amount/100)
			if err != nil {
				return err
			}
		} else {
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
		}

		err = s.orderRepo.FinishOrder(ctx, order.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
