package interfaces

import "github.com/streadway/amqp"

// AmpqClient интерфейс RabbitMQ клиента
type AmpqClient interface {
	Connect(dsn string) error
	Reconnect() error
	MigrateDurableQueues(queues ...string)
	Channel() *amqp.Channel
	Consume(queue string) (<-chan amqp.Delivery, error)
	Publish(queue string, msgBody []byte) error
	Stop()
}
