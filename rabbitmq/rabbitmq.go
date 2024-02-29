package rabbitmq

import (
	"context"
	"fmt"

	"github.com/pingcap/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	cfg  Config
	conn *amqp.Connection
	ch   *amqp.Channel
	q    []amqp.Queue
}

func NewRabbitMQ(options ...option) (*RabbitMQ, error) {
	rmq := &RabbitMQ{}

	for _, opt := range options {
		if err := opt(rmq); err != nil {
			return nil, err
		}
	}

	return rmq, nil
}

func (rmq *RabbitMQ) AddMigrates(migates ...Migrate) error {
	for _, migrate := range migates {
		if err := migrate(rmq); err != nil {
			return err
		}
	}

	return nil
}

func (rmq *RabbitMQ) connect() error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", rmq.cfg.Username, rmq.cfg.Password, rmq.cfg.Host, rmq.cfg.Port))
	if err != nil {
		return errors.Wrap(err, ErrConnectionFailed)
	}
	rmq.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, ErrCreateChannel)
	}
	rmq.ch = ch

	return nil
}

func (rmq *RabbitMQ) Start(ctx context.Context) error {
	if err := rmq.connect(); err != nil {
		return err
	}

	return nil
}

func (rmq *RabbitMQ) Stop(ctx context.Context) error {
	if err := rmq.ch.Close(); err != nil {
		return errors.Wrap(err, ErrCloseChannel)
	}

	if err := rmq.conn.Close(); err != nil {
		return errors.Wrap(err, ErrConnectionClose)
	}

	return nil
}
