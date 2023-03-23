package ampq

import (
	"github.com/streadway/amqp"

	"github.com/atrian/go-notify-customer/internal/interfaces"
)

type Client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     interfaces.Logger
	dsn        string
}

func New(dsn string, logger interfaces.Logger) *Client {
	client := Client{
		logger: logger,
		dsn:    dsn,
	}

	return &client
}

func NewWithConnection(dsn string, logger interfaces.Logger) *Client {
	client := Client{
		logger: logger,
		dsn:    dsn,
	}

	err := client.Connect(dsn)
	if err != nil {
		logger.Error("Can't connect AMPQ", err)
	}

	return &client
}

func (c *Client) Connect(dsn string) error {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		c.logger.Error("Can't connect AMPQ", err)
		return err
	}
	c.connection = conn

	channel, err := conn.Channel()
	if err != nil {
		c.logger.Error("Can't create AMPQ channel", err)
		return err
	}
	c.channel = channel

	return nil
}

// Reconnect переподключение к RabbitMQ
func (c *Client) Reconnect() error {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		c.logger.Error("Can't connect AMPQ", err)
		return err
	}
	c.connection = conn

	channel, err := conn.Channel()
	if err != nil {
		c.logger.Error("Can't create AMPQ channel", err)
		return err
	}
	c.channel = channel

	return nil
}

func (c *Client) Channel() *amqp.Channel {
	return c.channel
}

// MigrateDurableQueues создает Durable очереди в RabbitMQ
func (c *Client) MigrateDurableQueues(queues ...string) {
	for _, queue := range queues {
		if queue == "" {
			continue
		}

		_, err := c.Channel().QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			c.logger.Error("Can't declare queue", err)
		}
	}
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	consume, err := c.Channel().Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		c.logger.Error("AMPQ consume error", err)
	}

	return consume, nil
}

func (c *Client) Publish(queue string, msgBody []byte) error {
	err := c.Channel().Publish("", queue, false, false,
		amqp.Publishing{
			Headers:     nil,
			ContentType: "text/json",
			Body:        msgBody,
		})

	if err != nil {
		c.logger.Error("AMPQ publish error", err)
		return err
	}

	return nil
}

func (c *Client) Stop() {
	if c.channel != nil {
		chErr := c.channel.Close()
		if chErr != nil {
			c.logger.Error("Channel close error", chErr)
		}
	}

	if c.connection != nil {
		cErr := c.connection.Close()
		if cErr != nil {
			c.logger.Error("Connection close error", cErr)
		}
	}
}
