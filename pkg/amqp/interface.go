package amqp

import (
	"context"

	sdk "github.com/rabbitmq/amqp091-go"
)

type Amqp interface {
	PublishWithContext(_ context.Context, exchange, key string, mandatory, immediate bool, msg sdk.Publishing) error
	ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args sdk.Table) (<-chan sdk.Delivery, error)
}
