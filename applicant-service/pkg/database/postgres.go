package database

import (
	"fineDeedSystem/applicant-service/configs"
	"fineDeedSystem/applicant-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB(config configs.PostgresConfig) *gorm.DB {
	dsn := "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.DBName + " port=" + config.Port + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	// Auto-migrate the Applicant model
	db.AutoMigrate(&models.Applicant{})

	return db
}
