package databus

import (
	"context"
	"encoding/json"

	"github.com/Gunga-D/service-godzilla-soft-site/pkg/amqp"
	sdk_amqp "github.com/rabbitmq/amqp091-go"
)

type client struct {
	amqp amqp.Amqp
}

func NewClient(ctx context.Context) *client {
	return &client{
		amqp: amqp.Get(ctx, []string{
			"queue-item-out-of-stock",
		}),
	}
}

func (c *client) PublishDatabusItemOutOfStockDTO(ctx context.Context, msg ItemOutOfStockDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-item-out-of-stock", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}
