package databus

import (
	"context"

	sdk_amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ConsumeDatabusChangeItemState(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusNewUserEmail(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusSendToEmail(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusNewUserSteamLink(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusNeuroTask(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusNeuroNewItems(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusTelegramRegistration(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
}
