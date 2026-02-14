package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/TubagusAldiMY/go-template/internal/infrastructure/config"
	"github.com/TubagusAldiMY/go-template/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(cfg config.RabbitMQConfig) (*RabbitMQ, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	logger.Info("rabbitmq connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() error {
	var err error
	if r.channel != nil {
		if e := r.channel.Close(); e != nil {
			err = e
		}
	}
	if r.conn != nil {
		if e := r.conn.Close(); e != nil {
			err = e
		}
	}
	logger.Info("rabbitmq connection closed")
	return err
}

func (r *RabbitMQ) DeclareQueue(name string, durable, autoDelete bool) error {
	_, err := r.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (r *RabbitMQ) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	return r.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return r.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func (r *RabbitMQ) Consume(queueName, consumer string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,
		consumer,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQ) Health() error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	if r.channel == nil {
		return fmt.Errorf("rabbitmq channel is nil")
	}
	return nil
}
