package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Connection struct {
	*amqp.Connection
}

type Channel struct {
	*amqp.Channel
}

func ConnectToRabbitMQ(url string, retries int) (*Connection, *Channel, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < retries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return &Connection{conn}, &Channel{ch}, nil
}
