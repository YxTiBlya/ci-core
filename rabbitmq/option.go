package rabbitmq

type option = func(*RabbitMQ) error

func WithConfig(cfg Config) option {
	return func(rmq *RabbitMQ) error {
		rmq.cfg = cfg
		return nil
	}
}
