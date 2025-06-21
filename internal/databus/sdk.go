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
			"queue-new-user-email",
			"queue-send-to-email",
			"queue-new-user-steam-link",
			"queue-neuro-task",
			"queue-neuro-new-items",
			"queue-telegram-registration",
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

func (c *client) PublishDatabusNewUserEmail(ctx context.Context, msg NewUserEmailDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-new-user-email", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusSendToEmail(ctx context.Context, msg SendToEmailDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-send-to-email", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusNewUserSteamLink(ctx context.Context, msg NewUserSteamLinkDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-new-user-steam-link", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusNeuroTask(ctx context.Context, msg NeuroTaskDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-neuro-task", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusNeuroNewItems(ctx context.Context, msg NeuroNewItemsDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-neuro-new-items", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) PublishDatabusTelegramRegistration(ctx context.Context, msg TelegramRegistrationDTO) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.amqp.PublishWithContext(ctx, "", "queue-telegram-registration", true, false, sdk_amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}

func (c *client) ConsumeDatabusChangeItemState(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-item-change-item-state", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusNewUserEmail(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-new-user-email", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusSendToEmail(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-send-to-email", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusNewUserSteamLink(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-new-user-steam-link", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusNeuroTask(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-neuro-task", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusNeuroNewItems(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-neuro-new-items", "", false, false, false, false, nil)
}

func (c *client) ConsumeDatabusTelegramRegistration(ctx context.Context) (<-chan sdk_amqp.Delivery, error) {
	return c.amqp.ConsumeWithContext(ctx, "queue-telegram-registration", "", false, false, false, false, nil)
}
