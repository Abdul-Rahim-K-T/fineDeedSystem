package rabbitmq

import (
	"encoding/json"
	"fineDeedSystem/employer-service/internal/models"
	"fineDeedSystem/employer-service/internal/repository/postgres"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type Consumer struct {
	repo       *postgres.EmployerRepository
	rabbitMQCh *amqp.Channel
}

func NewConsumer(repo *postgres.EmployerRepository, rabitMQCh *amqp.Channel) *Consumer {
	return &Consumer{
		repo:       repo,
		rabbitMQCh: rabitMQCh,
	}
}

func (c *Consumer) Start() {
	log.Println("Starting RabbitMQ consumer...")
	msgs, err := c.rabbitMQCh.Consume(
		"employer_list_queue", // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   //args
	)

	if err != nil {
		log.Fatalf("Error setting up consumer: %v", err)
	}

	for d := range msgs {
		// Decode the message to determine the requested action
		var request map[string]string
		if err := json.Unmarshal(d.Body, &request); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue
		}

		// Check action type
		if action, ok := request["action"]; ok && action == "get_all_employers" {
			log.Printf("Processing action: %s", action)
			go c.handleGetAllEmployers(d.ReplyTo, d.CorrelationId)
		} else {
			log.Printf("Unknown action recieved: %s", request["action"])
		}
	}
}

func (c *Consumer) handleGetAllEmployers(replyTo string, correlationId string) {
	// // Ensure the replyTo queue exists
	// _, err := c.rabbitMQCh.QueueDeclare(
	// 	replyTo, // Name of the replyTo queue
	// 	false,   // Durable
	// 	false,   // Delete when unused
	// 	true,    // Exclusive
	// 	false,   // No-wait
	// 	nil,     // Arguments
	// )
	// if err != nil {
	// 	log.Printf("Error declaring replyTo queue '%s': %v", replyTo, err)
	// 	return
	// }

	log.Println("Fetching all employers from the database...")
	employers, err := c.repo.GetAllEmployers()
	if err != nil {
		log.Printf("Error fetching employers: %v", err)
		return
	}

	var wg sync.WaitGroup
	employerResponses := make([]models.EmployerResponse, len(employers))

	for i, employer := range employers {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int, employer models.Employer) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			// Log each employer's details
			log.Printf("ID: %d, Name: %s, Contact: %s", employer.ID, employer.Name, employer.Phone)

			// Convert to EmployerResponse
			employerResponses[i] = models.EmployerResponse{
				ID:          employer.ID,
				Name:        employer.Name,
				Email:       employer.Email,
				Phone:       employer.Phone,
				CompanyName: employer.CompanyName,
				CreatedAt:   employer.CreatedAt,
				UpdatedAt:   employer.UpdatedAt,
			}
		}(i, employer) // Pass loop variables to Goroutine
	}

	wg.Wait() // Wait for all goroutin to finis

	// Marshal the employer response data into JSON
	body, err := json.Marshal(employerResponses)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}

	// Publish the response back to the ReplyTo queue with the original CorrelationId
	log.Printf("Sending response to queue '%s'  with correlation ID: %s", replyTo, correlationId)
	err = c.rabbitMQCh.Publish(
		"",      // exchange
		replyTo, //routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: correlationId,
		})
	if err != nil {
		log.Printf("Error publishing responses: %v", err)
	} else {
		log.Println("Response successfully published")
	}
}
