package models

import "time"

type EmployerResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	CompanyName string    `json:"company_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e Employer) ToResponse() EmployerResponse {
	return EmployerResponse{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Phone:       e.Phone,
		CompanyName: e.CompanyName,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
