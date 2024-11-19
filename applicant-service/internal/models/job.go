package models

import "time"

type Job struct {
	ID          uint      `gorm:"primaryKey:autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CompanyID   uint      `json:"company_id"`
	Location    string    `json:"location"`
	Salary      float64   `json:"salary"`
	PostedAt    time.Time `json:"posted_at"`
	ClosedAt    time.Time `json:"closed_at"`
}
