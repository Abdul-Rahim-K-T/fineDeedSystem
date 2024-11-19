package repository

import "fineDeedSystem/applicant-service/internal/models"

type ApplicantRepository interface {
	Create(applicant *models.Applicant) error
	GetByEmail(email string) (*models.Applicant, error)
	GetApplicantByID(id int) (*models.Applicant, error)
	UpdateApplicant(applicant *models.Applicant) error
	DeleteApplicant(id uint) error
}
