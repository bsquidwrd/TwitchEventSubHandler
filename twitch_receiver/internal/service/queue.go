package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "twitch"

type queueService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func newQueueService() *queueService {
	service := &queueService{}
	service.connect()

	// Handle disconnects and attempt to connect again
	// unless the err is nil, indicating user shutdown/close
	go func() {
		for {
			errChan := make(chan *amqp.Error)
			service.ch.NotifyClose(errChan)
			err := <-errChan

			if err != nil {
				service.connect()
			} else {
				break
			}
		}
	}()

	slog.Debug("Queue connected successfully")

	return service
}

func (q *queueService) connect() {
	conn, err := amqp.Dial(os.Getenv("QUEUE_URL"))
	if err != nil {
		panic(err)
	}
	q.conn = conn

	ch, err := q.conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		panic(err)
	}
	q.ch = ch
}

func (q *queueService) cleanup() {
	defer q.conn.Close()
	defer q.ch.Close()
}

func (q *queueService) Ping() error {
	if q.ch.IsClosed() {
		return errors.New("queue is closed")
	}
	return nil
}

func (q *queueService) Publish(topic string, body interface{}) error {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.ch.PublishWithContext(ctx,
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJson,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
