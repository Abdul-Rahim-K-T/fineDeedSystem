package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Publisher struct for publishing messages to RabbitMQ
type Publisher struct {
	Channel *amqp.Channel
	Queue   string
}

// NewPublisher creates a new publisher
func NewPublisher(Channel *amqp.Channel, queue string) *Publisher {
	return &Publisher{
		Channel: Channel,
		Queue:   queue,
	}
}

// PublishMessage publishes a message to the RabbitMQ queue
func (p *Publisher) PublishMessage(body []byte) error {
	err := p.Channel.Publish(
		"",      // exchange
		p.Queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return nil
	}
	log.Printf("Message published to queue '%s'", p.Queue)
	return nil
}

// // Consumer struct for consuming messages from RabbitMQ
// type Consumer struct {
// 	Channel *amqp.Channel
// 	Queue   string
// }

// // NewConsumer creates a new Consumer
// func NewConsumer(channel *amqp.Channel, queue string) *Consumer {
// 	return &Consumer{
// 		Channel: channel,
// 		Queue: queue,
// 	}
// }

// // ConsumeMessages consumes messages from the RabbitMQ queue
// func (c *Consumer) ConsumeMessages(handler func(amqp.Delivery)) error {
// 	msgs, err := c.Channel.Consume(
// 		c.Queue,   // queue
// 		"",     // Consumer
// 		true,    // auto-ack
// 		false, // exclusive
// 		false,  // no-local
// 		false,   // no-wait
// 		nil,   // args
// 	)
// 	if err != nil {
// 		log.Printf("Error setting up consumer: %v", err)
// 		return err
// 	}

// 	go func() {
// 		for d := range msgs {
// 			handler(d)
// 		}
// 	}()
// 	log.Printf("Consumer started for queue '%s'", c.Queue)
// 	return nil
// }
