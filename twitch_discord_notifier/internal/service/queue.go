package service

import (
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

func (q *queueService) StartConsuming(queueName string, topics []string, callback func(amqp.Delivery)) {
	queue, err := q.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	for _, topic := range topics {
		err = q.ch.QueueBind(
			queue.Name,   // queue name
			topic,        // routing key
			exchangeName, // exchange
			false,        // no wait
			nil,          // args
		)
		if err != nil {
			panic(err)
		}

		msgs, err := q.ch.Consume(
			queue.Name, // queue
			"",         // consumer
			true,       // auto ack
			false,      // exclusive
			false,      // no local
			false,      // no wait
			nil,        // args
		)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range msgs {
				m := msg
				go func() {
					startTime := time.Now().UTC()
					callback(m)
					endTime := time.Now().UTC()
					slog.Debug(" [x] Processed message", "topic", m.RoutingKey, "millisecondstoprocess", endTime.Sub(startTime).Milliseconds())
				}()
			}
		}()
	}

	slog.Debug("Successfully started consuming from queue", "queuename", queue.Name)
}
