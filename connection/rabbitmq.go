package connection

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessageBody struct {
	Data      []byte
	Type      string
	RequestID string
}

type RabbitMQMessage struct {
	Queue       string
	ContentType string
	Body        RabbitMQMessageBody
}

type RabbitMQConnection struct {
	name      string
	conn      *amqp.Connection
	channel   *amqp.Channel
	exchange  string
	queue     string
	err       chan error
	connected chan struct{}
	config    config.RabbitMQConfig
	logger    *logger.Logger
}

func NewRabbitMQConnection(config config.RabbitMQConfig, logger *logger.Logger, connectionName, exchange, queue string) *RabbitMQConnection {
	return &RabbitMQConnection{
		name:      connectionName,
		exchange:  exchange,
		queue:     queue,
		err:       make(chan error),
		connected: make(chan struct{}),
		config:    config,
		logger:    logger,
	}
}

func (c *RabbitMQConnection) Connect() error {
	var err error

	c.logger.Info(
		context.TODO(),
		"Connecting to RabbitMQ",
		"name", c.name,
		"queue", c.queue,
		"exchange", c.exchange,
	)

	host := net.JoinHostPort(c.config.Host, strconv.Itoa(c.config.Port)) + "/" + c.config.VHost
	url := fmt.Sprintf("amqp://%s:%s@%s", c.config.User, c.config.Password, host)
	c.conn, err = amqp.DialConfig(url, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 20*time.Second)
		},
	})
	if err != nil {
		c.logger.Error(
			context.TODO(),
			"error in creating rabbitmq connection",
			"error", err,
		)
		return err
	}

	go func() {
		amqpErr := <-c.conn.NotifyClose(make(chan *amqp.Error))
		if amqpErr != nil {
			c.logger.Error(
				context.TODO(),
				"connection Closed",
				"error", amqpErr.Error(),
			)
			c.err <- amqpErr
		}
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		c.logger.Error(
			context.TODO(),
			"error in creating rabbitmq channel",
			"error", err,
		)
		return err
	}
	if err = c.channel.ExchangeDeclare(
		c.exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		c.logger.Error(
			context.TODO(),
			"error in Exchange Declare",
			"error", err,
		)
		return err
	}

	err = c.BindQueue()
	if err != nil {
		c.logger.Error(
			context.TODO(),
			"error in binding queue",
			"error", err,
		)
		return err
	}

	go func() {
		c.connected <- struct{}{}
	}()

	c.logger.Info(
		context.TODO(),
		"Connected to RabbitMQ",
		"name", c.name,
	)

	return nil
}

func (c *RabbitMQConnection) ConnectionOpener() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
			case err := <-c.err:
				if err == nil {
					continue
				}
				c.logger.Error(
					context.TODO(),
					"Error in RabbitMQ connection",
					"error", err,
				)

				err = c.Close()
				if err != nil {
					c.logger.Error(
						context.TODO(),
						"error in closing",
						"error", err,
					)
				}
				err = c.Connect()
				if err != nil {
					c.logger.Error(
						context.TODO(),
						"error in reconnecting",
						"error", err,
					)
					time.Sleep(5 * time.Second)
					go func() {
						c.err <- errors.New("error in reconnecting")
					}()
				}
			}
		}
	}()
}

func (c *RabbitMQConnection) Consume(prefetchCount int) (<-chan amqp.Delivery, error) {
	err := c.channel.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		c.logger.Error(
			context.TODO(),
			"Error setting qos",
			"error", err,
		)
	}

	deliveries, err := c.channel.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (c *RabbitMQConnection) HandleConsumedDeliveries(ctx context.Context, prefetchCount int, fn func(context.Context, *logger.Logger, RabbitMQConnection, amqp.Delivery)) {
	workerPoolSize := runtime.NumCPU()
	workerPool := make(chan struct{}, workerPoolSize)
	for {
		select {
		case <-ctx.Done():
			c.logger.Error(
				context.TODO(),
				"Error consuming queue",
				"error", ctx.Err(),
			)
			return
		case <-c.connected:
			c.logger.Info(
				context.TODO(),
				"Reconnected to RabbitMQ",
			)
			deliveries, err := c.Consume(prefetchCount)
			if err != nil {
				c.logger.Error(
					context.TODO(),
					"Error consuming queue",
					"error", err,
				)
				continue
			}
			c.processDeliveriesWithWorkerPool(ctx, deliveries, fn, workerPool)
		}
	}
}

func (c *RabbitMQConnection) processDeliveriesWithWorkerPool(ctx context.Context, deliveries <-chan amqp.Delivery, fn func(context.Context, *logger.Logger, RabbitMQConnection, amqp.Delivery), workerPool chan struct{}) {
	for delivery := range deliveries {
		workerPool <- struct{}{} // Acquire a worker goroutine
		go func(delivery amqp.Delivery) {
			defer func() {
				if r := recover(); r != nil {
					c.logger.Error(
						context.TODO(),
						"Error in processing message",
						"panic", r,
					)

					err := delivery.Reject(false)
					if err != nil {
						c.logger.Error(
							context.TODO(),
							"failed to Reject",
							"error", err,
						)
					} else {
						c.logger.Info(
							context.TODO(),
							"reject message after we faced panic",
						)
					}
				}
				<-workerPool // Release the worker goroutine
			}()

			fn(ctx, c.logger, *c, delivery)

		}(delivery)
	}
}

func (c *RabbitMQConnection) Publish(ctx context.Context, message RabbitMQMessage) error {

	if c == nil || c.conn == nil {
		return errors.New("rabbitmq connection is nil")
	}

	if c.conn.IsClosed() {
		return errors.New("rabbitmq connection is closed")
	}

	publishing := amqp.Publishing{
		Headers:     amqp.Table{"type": message.Body.Type},
		ContentType: message.ContentType,
		Body:        message.Body.Data,
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := c.channel.PublishWithContext(ctx, "", message.Queue, false, false, publishing); err != nil {
		c.logger.Error(
			context.TODO(),
			"error in Publishing",
			"error", err,
		)
		return err
	}

	return nil
}

func (c *RabbitMQConnection) BindQueue() error {
	if _, err := c.channel.QueueDeclare(c.queue, true, false, false, false, nil); err != nil {
		c.logger.Error(
			context.TODO(),
			"error in declaring the queue",
			"error", err,
		)
		return err
	}
	if err := c.channel.QueueBind(c.queue, "", c.exchange, false, nil); err != nil {
		c.logger.Error(
			context.TODO(),
			"error in binding the queue",
			"error", err,
		)
		return err
	}
	return nil
}

func (c *RabbitMQConnection) Close() error {
	if c.conn != nil {
		c.logger.Info(
			context.TODO(),
			"Closing the rabbitmq connection",
		)
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	if c.channel != nil {
		c.logger.Info(
			context.TODO(),
			"Closing the rabbitmq channel",
		)
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	return nil
}
