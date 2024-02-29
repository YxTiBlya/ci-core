package rabbitmq

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Migrate = func(*RabbitMQ) error

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func WithQueue(cfg *QueueConfig) Migrate {
	return func(rm *RabbitMQ) error {
		q, err := rm.ch.QueueDeclare(
			cfg.Name,
			cfg.Durable,
			cfg.AutoDelete,
			cfg.Exclusive,
			cfg.NoWait,
			cfg.Args,
		)
		if err != nil {
			return errors.Wrap(err, ErrQueueDeclare)
		}

		if rm.q == nil {
			rm.q = make([]amqp.Queue, 0, 5) // TODO: to cfg?
		}
		rm.q = append(rm.q, q)

		return nil
	}
}
