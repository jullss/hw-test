package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	qName string
}

func NewRabbitClient(url, queueName string) (*RabbitClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitClient{
		conn:  conn,
		ch:    ch,
		qName: queueName,
	}, nil
}

func (r *RabbitClient) Publish(ctx context.Context, msg Notification) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	return r.ch.PublishWithContext(ctx,
		"",
		r.qName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (r *RabbitClient) Consume(ctx context.Context) (<-chan Notification, error) {
	deliveries, err := r.ch.Consume(
		r.qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	out := make(chan Notification)

	go func() {
		defer close(out)
		for d := range deliveries {
			var msg Notification
			if err := json.Unmarshal(d.Body, &msg); err == nil {
				out <- msg
			}
		}
	}()

	return out, nil
}

func (r *RabbitClient) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
