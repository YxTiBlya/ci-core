package rabbitmq

import (
	"github.com/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rmq *RabbitMQ) Consume(qName, consumer string, ack, excl, nlocal, nwait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	msgs, err := rmq.ch.Consume(
		qName,
		consumer,
		ack,
		excl,
		nlocal,
		nwait,
		args,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrConsume)
	}

	return msgs, nil
}
