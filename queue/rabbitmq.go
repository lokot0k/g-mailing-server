package queue

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func NewRabbitMQ(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to RabbitMQ after 10 attempts: %w", err)
}

func Publish(conn *amqp.Connection, queueName string, body []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
