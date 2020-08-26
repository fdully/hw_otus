package rabbit

import (
	"context"
	"io"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(body []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context) <-chan amqp.Delivery
}

type AMQP interface {
	Publisher
	Subscriber
	io.Closer
}

type connector struct {
	exchange string
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    amqp.Queue
	QOS      int
	msgCh    <-chan amqp.Delivery
}

func NewConnector(ctx context.Context, url, exchange, queueName string, qos int) AMQP {
	c := &connector{exchange: exchange, QOS: qos}
	c.setupConn(ctx, url)
	c.declareQueue(ctx, queueName)

	return c
}

func (c *connector) Publish(body []byte) error {
	err := c.channel.Publish(c.exchange, c.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return err
}

func (c *connector) Subscribe(ctx context.Context) <-chan amqp.Delivery {
	c.setupQOS(ctx, c.QOS)
	c.setupMsgCh(ctx)

	return c.msgCh
}

func (c *connector) setupConn(ctx context.Context, url string) {
	logger := logging.FromContext(ctx)

	var err error
	c.conn, err = amqp.Dial(url)
	if err != nil {
		logger.Errorf("Can't connect to AMQP: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		logger.Errorf("Can't create a amqpChannel: %w", err)
	}
}

func (c *connector) declareQueue(ctx context.Context, name string) {
	var err error
	c.queue, err = c.channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		logger := logging.FromContext(ctx)
		logger.Errorf("Could not declare %s queue: %w", name, err)
	}
}

func (c *connector) setupQOS(ctx context.Context, count int) {
	if err := c.channel.Qos(count, 0, false); err != nil {
		logger := logging.FromContext(ctx)
		logger.Errorf("Could not configure QoS: %w", err)
	}
}

func (c *connector) setupMsgCh(ctx context.Context) {
	var err error
	c.msgCh, err = c.channel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger := logging.FromContext(ctx)
		logger.Errorf("Could not register consumer: %w", err)
	}
}

func (c *connector) Close() error {
	return c.conn.Close()
}
