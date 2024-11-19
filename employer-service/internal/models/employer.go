package models

import (
	"time"

	"gorm.io/gorm"
)

type Employer struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	CompanyName string         `json:"company_name"`
	Password    string         `json:"password"`              // Exclude this field
	Jobs        []Job          `gorm:"foreignKey:EmployerID"` // One-to-many relationship
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
