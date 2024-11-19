package main

import (
	"fineDeedSystem/applicant-service/configs"
	appHttp "fineDeedSystem/applicant-service/internal/delivery/http"
	repository "fineDeedSystem/applicant-service/internal/repository/postgres"
	"fineDeedSystem/applicant-service/internal/usecase"
	"fineDeedSystem/applicant-service/pkg/database"
	"fineDeedSystem/applicant-service/pkg/redisclient"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func initDatabase() *gorm.DB {
	dbConfig := configs.GetPostgresConfig()
	return database.InitPostgresDB(dbConfig)
}

func initRedis() *redisclient.RedisClient {
	redisConfig := configs.GetRedisConfig()
	return redisclient.NewRedisClient(redisConfig.Host, redisConfig.Port, redisConfig.Password)
}

// func initRedis() *database.RedisClient {
// 	redisConfig := configs.GetRedisConfig()
// 	return database.InitRedis(redisConfig)
// }

// func initRabbitMQ() (*rabbitmq.Connection, *rabbitmq.Channel, amqp.Queue) {
// 	rabbitMQURL := os.Getenv("RABBITMQ_URL")
// 	conn, ch, err := rabbitmq.ConnectToRabbitMQ(rabbitMQURL, 5)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}

// 	q, err := ch.QueueDeclare(
// 		"applicant_update_queue",	// name
// 		false,						// durable
// 		false,						// delete when unused
// 		false,						// exclusive
// 		false,						// no-wait
// 		nil,						// arguments
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to declare a queue: %v", err)
// 	}
// 	return conn, ch, q
// }

// func startRabbitMQConsumer(repo *postgres.ApplicantRepository, ch *rabbitmq.Channel) {
// 	log.Println("Starting RabbiMQ Consumer")
// 	consumer := rabbitmq.NewConsumer(repo, ch.Channel)
// 	go consumer.Start()
// }

// func startHTTPServer(port string, router *mux.Router) *http.Server {
// 	srv := &http.Server{Addr: ":" + port, Handler: router}
// 	go func() {
// 		log.Print("HTTP server running on port %s", port)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("ListenAndServe(): %v", err)
// 		}
// 	}()
// 	return srv
// }

// func StartGRPCServer(port string, applicantUsecase *usecase.ApplicantUsecase) *grpc.Server {
// 	listener, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		log.Fatalf("Failed to listen on port %s: %v", port, err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	applicantHandler := grpcHandler.NewApplicantGrpcHandler(applicantUsecase)
// 	shared.RegistereApplicantServiceServer(grpcServer, applicantHandler)
// 	go func() {
// 		log.Printf("gRPC server running on port %s", port)
// 		if err := grpcServer.Serve(listener); err != nil {
// 			log.Fatalf("Failed to serve gRPC server: %v", err)
// 		}
// 	}()
// 	return grpcServer
// }

// func waitForShutdown(srv *http.Server, grpcServer *grpc.Server) {
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)
// 	<-quit

// 	log.Println("Shutting down the server...")
// 	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatal("Server forced to shutdown:", err)
// 	}
// 	grpcServer.GracefulStop()
// 	log.Println("Server exiting")
// }

func main() {
	log.Println("Starting applicant-service...")

	initEnv()

	db := initDatabase()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Error getting DB from DB: %v", err)
		}
		sqlDB.Close()
	}()

	redisclient := initRedis()
	applicantRepo := repository.NewApplicantRepository(db)
	applicantUsecase := usecase.NewApplicantUsecase(applicantRepo, redisclient)

	router := mux.NewRouter()
	appHttp.RegisterApplicantRoutes(router, applicantUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
