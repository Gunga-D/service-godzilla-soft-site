package databus

import (
	"context"

	sdk_amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ConsumeDatabusChangeItemState(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
	ConsumeDatabusQuickUserRegistration(ctx context.Context) (<-chan sdk_amqp.Delivery, error)
}
