package database

import (
	"fmt"
	"log"

	"fineDeedSystem/employer-service/configs"
	"fineDeedSystem/employer-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgresDB initializes the PostgreSQL database and performs auto-migration
func InitPostgresDB(config configs.PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Perform auto-migration for models
	err = db.AutoMigrate(&models.Employer{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	return db
}
