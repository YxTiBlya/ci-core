package rabbitmq

var (
	ErrConnectionFailed = "connection failed"
	ErrConnectionClose  = "connection close is failed"

	ErrCreateChannel = "create channel is failed"
	ErrCloseChannel  = "channel closing is failed"

	ErrQueueDeclare = "queue declare is failed"

	ErrPublish = "publish is failed"

	ErrConsume = "consume is failed"
)
