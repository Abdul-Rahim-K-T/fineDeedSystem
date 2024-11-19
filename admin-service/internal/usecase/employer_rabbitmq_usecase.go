package usecase

import (
	"encoding/json"
	"errors"
	"fineDeedSystem/admin-service/internal/models"
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type EmployerRabbitMQUsecase struct {
	rabbitMQCh *amqp.Channel
}

func NewEmployerRabbitMQUsecase(rabbitMQCh *amqp.Channel) *EmployerRabbitMQUsecase {
	return &EmployerRabbitMQUsecase{
		rabbitMQCh: rabbitMQCh,
	}
}

func (u *EmployerRabbitMQUsecase) GetAllEmployers() ([]*models.Employer, error) {
	log.Println("GetAllEmployers via RabbitMQ /employer_rabbitmq_usecase/")
	request := map[string]string{"action": "get_all_employers"}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseQueue, err := u.rabbitMQCh.QueueDeclare(
		"",    // name (empty means generate a unique name)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Error declaring response queue:", err)
		return nil, err
	}

	correlationId := strconv.FormatInt(time.Now().UnixNano(), 10)
	err = u.rabbitMQCh.Publish(
		"",                    // exchange
		"employer_list_queue", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       responseQueue.Name,
			CorrelationId: correlationId,
		})
	if err != nil {
		log.Println("Error publishing message: ", err)
		return nil, err
	}

	msgs, err := u.rabbitMQCh.Consume(
		responseQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Println("Error consuming response queue:", err)
		return nil, err
	}

	responseChan := make(chan []*models.Employer)

	go func() {
		for d := range msgs {
			if d.CorrelationId == correlationId {
				var employers []*models.Employer
				if err := json.Unmarshal(d.Body, &employers); err != nil {
					log.Println("Received matching response with CorrelationId:", correlationId)
					log.Println("Error unmarshalling response:", err)
					responseChan <- nil
					return
				}
				responseChan <- employers
				return
			}
		}
	}()

	select {
	case employers := <-responseChan:
		if employers == nil {
			return nil, errors.New("received invalid response")
		}
		return employers, nil
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for response from RabbitMQ")
		return nil, errors.New("timeout waiting for response from RabbitMQ")
	}
}
