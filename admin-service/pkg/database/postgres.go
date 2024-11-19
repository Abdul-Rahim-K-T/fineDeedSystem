package database

import (
	"fineDeedSystem/admin-service/configs"
	"fineDeedSystem/admin-service/internal/models"
	"fmt"

	// "fineDeedSystem/admin-service/internal/repository/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB(config configs.PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Auto-migrate the Admin model
	err = db.AutoMigrate(&models.Admin{})
	if err != nil {
		panic("failed to migrate database")
	}
	return db
}
