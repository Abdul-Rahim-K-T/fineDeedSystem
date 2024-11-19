package main

import (
	"context"

	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fineDeedSystem/employer-service/configs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	"gorm.io/gorm"

	httpRout "fineDeedSystem/employer-service/internal/delivery/http"

	"fineDeedSystem/employer-service/internal/rabbitmq"
	"fineDeedSystem/employer-service/internal/repository/postgres"
	"fineDeedSystem/employer-service/internal/usecase"
	"fineDeedSystem/employer-service/proto/fineDeedSystem/proto/shared"

	grpcHandler "fineDeedSystem/employer-service/internal/grpc"
	"fineDeedSystem/employer-service/pkg/database"

	"google.golang.org/grpc"
)

// Initialize environment variables
func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Initialize PostgreSQL database
func initDatabase() *gorm.DB {
	dbConfig := configs.GetPostgresConfig()
	return database.InitPostgresDB(dbConfig)
}

// Initialize Redis
func initRedis() *database.RedisClient {
	redisConfig := configs.GetRedisConfig()
	return database.InitRedis(redisConfig)
}

// Initialize RabbitMQ
func initRabbitMQ() (*rabbitmq.Connection, *rabbitmq.Channel, amqp.Queue) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	conn, ch, err := rabbitmq.ConnectToRabbitMQ(rabbitMQURL, 5)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	_, err = ch.QueueDeclare(
		"employer_list_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	q, err := ch.QueueDeclare(
		"employer_update_queue", // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a employer_update_queue: %v", err)
	}
	return conn, ch, q
}

// Start RabbitMQ consumer
func startRabbitMQConsumer(repo *postgres.EmployerRepository, ch *rabbitmq.Channel) {
	log.Println("Starting RabbitMQ Consumer")
	consumer := rabbitmq.NewConsumer(repo, ch.Channel)
	go consumer.Start()
}

// Start HTTP server
func startHTTPServer(port string, router *mux.Router) *http.Server {
	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("HTTP server running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}

// Start gRPC server
func StartGRPCServer(port string, employerUsecase *usecase.EmployerUsecase) *grpc.Server {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	employerHandler := grpcHandler.NewEmployerGrpcHandler(employerUsecase)
	shared.RegisterEmployerServiceServer(grpcServer, employerHandler)
	go func() {
		log.Printf("gRPC server running on port %s", port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	return grpcServer
}

func waitForShutdown(srv *http.Server, grpcServer *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Sutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	grpcServer.GracefulStop()
	log.Println()
	log.Println("Server exiting")
}

func main() {
	log.Println("Starting employer-service..")

	// Initialize environment variables
	initEnv()

	// Initialize PostgreSQL
	db := initDatabase()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Error getting DB from DB: %v", err)
		}
		sqlDB.Close()
	}()

	// Initialize Redis
	redisClient := initRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("Error closing Redis client: %v", err)
		}
	}()

	// Initialize RabbitMQ
	rabbitMQConn, rabbitMQChannel, q := initRabbitMQ()
	defer rabbitMQConn.Close()
	defer rabbitMQChannel.Close()

	// Initialize usecase and repositories
	employerRepo := postgres.NewEmployerRepository(db)
	employerUsecase := usecase.NewEmployerUsecase(employerRepo, redisClient, q.Name)

	// Start RabbitMQ Consumer
	startRabbitMQConsumer(employerRepo, rabbitMQChannel)

	// Start HTTP and gRPC servers
	router := mux.NewRouter()
	httpRout.RegisterEmployerRoutes(router, employerUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	grpcPort := "50051"

	// Start gRPC server in a goroutine and assign it to a variable
	var grpcServer *grpc.Server
	go func() {
		grpcServer = StartGRPCServer(grpcPort, employerUsecase)
	}()

	// Start HTTP
	httpServer := startHTTPServer(port, router)

	// Wait for interrupt signal to gracefully shutdown the server
	waitForShutdown(httpServer, grpcServer)

}
