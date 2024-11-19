package repository

import (
	"fineDeedSystem/applicant-service/internal/models"

	"gorm.io/gorm"
)

type ApplicantRepository struct {
	db *gorm.DB
}

func NewApplicantRepository(db *gorm.DB) *ApplicantRepository {
	return &ApplicantRepository{
		db: db,
	}
}

func (r *ApplicantRepository) Create(applicant *models.Applicant) error {
	return r.db.Create(applicant).Error
}

func (r *ApplicantRepository) GetByEmail(email string) (*models.Applicant, error) {
	var applicant models.Applicant
	if err := r.db.Where("email = ?", email).First(&applicant).Error; err != nil {
		return nil, err
	}
	return &applicant, nil
}

func (r *ApplicantRepository) GetApplicantByID(id int) (*models.Applicant, error) {
	var applicant models.Applicant
	err := r.db.Where("id = ?", id).First(&applicant).Error
	if err != nil {
		return nil, err
	}
	return &applicant, nil
}

func (r *ApplicantRepository) UpdateApplicant(applicant *models.Applicant) error {
	return r.db.Save(applicant).Error
}

func (r *ApplicantRepository) DeleteApplicant(id uint) error {
	return r.db.Delete(&models.Applicant{}, id).Error
}
