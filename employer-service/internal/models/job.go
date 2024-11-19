package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	Title       string   `json:"title"`
	Description string   `json:"description"`
	EmployerID  uint     `json:"employer_id"`
	Employer    Employer `gorm:"foreignKey:EmployerID"`
}
