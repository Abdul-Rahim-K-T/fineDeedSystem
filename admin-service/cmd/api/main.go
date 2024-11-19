package main

import (
	// "fineDeedSystem/admin-service/internal/grpc"
	"context"
	"fineDeedSystem/admin-service/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fineDeedSystem/admin-service/configs"

	"fineDeedSystem/admin-service/pkg/database"

	pb "fineDeedSystem/admin-service/proto/fineDeedSystem/proto/shared"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	httpRout "fineDeedSystem/admin-service/internal/delivery/http"
	"fineDeedSystem/admin-service/internal/repository/postgres"

	"google.golang.org/grpc"
)

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

func consumeRabbitMQMessages(ch *amqp.Channel, adminUsecase *usecase.AdminUsecase) {
	msgs, err := ch.Consume(
		"employer_update_queue", // queue
		"",                      // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			// Process the message
			log.Printf("Received a message: %s", msg.Body)
			adminUsecase.InvalidateCache(msg.Body)
		}
	}()
}

// func handleEmployerUpdate(ch *amqp.Channel, redisClient *database.RedisClient) {
// 	msgs, err := ch.Consume(
// 		"employer_update_queue", // queue
// 		"",                      // consumer
// 		true,                    // auto-ack
// 		false,                   // exlusive
// 		false,                   // no-local
// 		false,                   // no-wait
// 		nil,                     // args
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register a consumer: %v", err)
// 	}

// 	forever := make(chan bool)

// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)

// 			// Invalidate the cache here
// 			err := redisClient.InvalidateCache("employer_cache_key") // Replace with actual cache key
// 			if err != nil {
// 				log.Printf("Failed to invalidate cache: %v", err)
// 			} else {
// 				log.Println("Cache invalidated successfully")
// 			}
// 		}
// 	}()

// 	log.Println("Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Redis
	redisConfig := configs.GetRedisConfig()
	redisClient := database.InitRedis(redisConfig)

	// Initialize PostgreSQL
	db := database.InitPostgresDB(configs.GetPostgresConfig())
	adminRepo := postgres.NewAdminRepository(db)

	// Create a gRPC client
	grpcConn, err := grpc.Dial(os.Getenv("EMPLOYER_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to employer service: %v", err)
	}
	defer grpcConn.Close()

	grpcClient := pb.NewEmployerServiceClient(grpcConn)
	employerGrpcUsecase := usecase.NewEmployerGrpcUsecase(grpcClient)

	// Initialize RabbitMQ
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"employer_list_queue", // name
	// 	false,                 // durable
	// 	false,                 // delete when unused
	// 	false,                 // exclusive
	// 	false,                 // no-wait
	// 	nil,                   // arguments
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to declare a queue: %v", err)
	// }

	employerRabbitMQUsecase := usecase.NewEmployerRabbitMQUsecase(ch)

	// Create an instance of AdminUsecase with the gRPC client
	adminUsecase := usecase.NewAdminUsecase(adminRepo, redisClient, employerRabbitMQUsecase, employerGrpcUsecase)

	// Consume RabbitMQ messages
	consumeRabbitMQMessages(ch, adminUsecase)

	// Setup routes
	r := mux.NewRouter()
	httpRout.RegisterAdminRoutes(r, adminUsecase, grpcClient, redisClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Sever running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// // Start listening to the employer update queue
	// go func() {
	// 	msgs, err := ch.Consume(
	// 		"employer_update_queue",  // queue
	// 		"",          // Consumer
	// 		true,           // auto-ack
	// 		false,             // exclusive
	// 		false,             // no-local
	// 		false,              // no-wait
	// 		nil,                // args
	// 	)
	// 	if err != nil {
	// 		log.Fatalf("Failed to register a consumer: %v", err)
	// 	}

	// 	for msg := range msgs {
	// 		log.Printf("Received a message: %s", msg.Body)
	// 		// Invalidate cache
	// 		err := adminUsecase.InvalidateCache(context.Background())
	// 		if err != nil {
	// 			log.Printf("Failed to invalidate cache: %v", err)
	// 		}
	// 	}
	// }()

	// // Start the RabbitMQ consumer to handle employer updates
	// go handleEmployerUpdate(ch, redisClient)

	// Graceful shutdown
	gracefulShutdown(server)

}
