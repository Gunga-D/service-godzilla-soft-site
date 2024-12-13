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
			"queue-item-change-item-state",
			"queue-quick-user-registration",
		}),
	}
}

func (c *client) PublishDatabusChangeItemState(ctx context.Context, msg ChangeItemStateDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-item-change-item-state", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusQuickUserRegistration(ctx context.Context, msg QuickUserRegistrationDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-quick-user-registration", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) ConsumeDatabusChangeItemState(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-item-change-item-state", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusQuickUserRegistration(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-quick-user-registration", "", false, false, false, false, nil)
}
