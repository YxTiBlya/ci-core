package rabbitmq

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (rmq *RabbitMQ) Publish(
	ctx context.Context,
	exchange, key string,
	mandatory bool,
	immediate bool,
	msg *amqp.Publishing,
) error {
	if err := rmq.ch.PublishWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		*msg,
	); err != nil {
		return errors.Wrap(err, ErrPublish)
	}

	return nil
}
