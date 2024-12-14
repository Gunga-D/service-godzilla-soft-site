package send_to_email

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type handler struct {
	consumer         databus.Consumer
	yandexMailClient yandex_mail.Client
}

func NewHandler(consumer databus.Consumer, yandexMailClient yandex_mail.Client) *handler {
	return &handler{
		consumer:         consumer,
		yandexMailClient: yandexMailClient,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusSendToEmail(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus change item state: %v", err)
	}
	for msg := range msgs {
		var data databus.SendToEmailDTO
		json.Unmarshal(msg.Body, &data)

		log.Printf("[info] trying send to email %s %s\n", data.Email, data.Subject)

		err := h.yandexMailClient.SendMail([]string{
			data.Email,
		}, data.Subject, data.Body)
		if err != nil {
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
}
