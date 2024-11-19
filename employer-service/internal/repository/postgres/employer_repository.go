package postgres

import (
	"context"
	"errors"
	"fineDeedSystem/employer-service/internal/models"
	"fineDeedSystem/employer-service/internal/utils"
	"strconv"

	"gorm.io/gorm"
)

// EmployerRepository defines the methods for employer data operations.
type EmployerRepository struct {
	db *gorm.DB
}

// New EmployerRepository creates a new instance of EmployerRepository.
func NewEmployerRepository(db *gorm.DB) *EmployerRepository {
	return &EmployerRepository{db: db}
}

// CreateEmployer adds a new employer to the database.
func (r *EmployerRepository) CreateEmployer(employer *models.Employer) error {
	return r.db.Create(employer).Error
}

// GetEmployerByID retrieves an employer by their ID.
func (r *EmployerRepository) GetEmployerByID(id uint) (*models.Employer, error) {
	var employer models.Employer
	if err := r.db.First(&employer, id).Error; err != nil {
		return nil, err
	}
	return &employer, nil
}

// DeleteEmployer removes an employer from the database.
func (r *EmployerRepository) DeleteEmployer(id uint) error {
	return r.db.Delete(&models.Employer{}, id).Error
}

// ListEmployers retriives all employers from the database.
func (r *EmployerRepository) ListEmployers() ([]models.Employer, error) {
	var employer []models.Employer
	if err := r.db.Find(&employer).Error; err != nil {
		return nil, err
	}
	return employer, nil
}

func (r *EmployerRepository) GetEmployerByUsername(username string) (*models.Employer, error) {
	var employer models.Employer
	result := r.db.Where("username = ?", username).First(&employer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &employer, nil
}

func (r *EmployerRepository) GetByEmail(email string) (*models.Employer, error) {
	var employer models.Employer
	if err := r.db.Where("email = ?", email).First(&employer).Error; err != nil {
		return nil, err
	}
	return &employer, nil
}

// FindByEmail retrieves an employer by their email
func (r *EmployerRepository) FindByEmail(email string) (*models.Employer, error) {
	var employer models.Employer
	if err := r.db.Where("email = ?", email).First(&employer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No employer found with this email
		}
		return nil, err // An error occurred during the query
	}
	return &employer, nil // Employer found
}

func (r *EmployerRepository) GetEmployerProfile(employerID string) (*models.Employer, error) {
	var employer models.Employer
	if err := r.db.Where("id = ?", employerID).First(&employer).Error; err != nil {
		return nil, err
	}
	return &employer, nil
}

// UpdateEmployerProfile updates an employer's profile.
func (r *EmployerRepository) UpdateEmployerProfile(employerID string, updatedProfile *models.Employer) error {
	return r.db.Model(&models.Employer{}).Where("id = ?", employerID).Updates(updatedProfile).Error
}

// ChangePassword changes the password of an employer
func (r EmployerRepository) ChangePassword(employerID, oldPassword, newPassword string) error {
	var employer models.Employer
	err := r.db.Where("id = ?", employerID).First(&employer).Error
	if err != nil {
		return err
	}

	// Check if the old  password is correct (assuming a HashPassword function)
	if !utils.CheckPasswordHash(oldPassword, employer.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash the new password before saving (assuming a HashPassword function)
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the password
	employer.Password = hashedPassword
	return r.db.Save(&employer).Error

}

// GetEmployerDashboard retrieves data for the employer's dashboard.
func (r *EmployerRepository) GetEmployerDashboard(employerID string) (interface{}, error) {
	// Logic for gathering dashboard data, such as job statistics, etc.
	// Placeholder return
	return nil, nil
}

// ListPostedJobs retrieves all jobs posted by an employer.
func (r *EmployerRepository) ListPostedJobs(employerID string) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Where("employer_id = ?", employerID).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// CreateJobPosting creates a new job posting.
func (r *EmployerRepository) CreateJobPosting(employerID string, job *models.Job) error {
	// Convert employerID from string to uint
	id, err := strconv.ParseUint(employerID, 10, 32)
	if err != nil {
		return errors.New("invalid employer ID")
	}
	job.EmployerID = uint(id) // Set the employer ID
	return r.db.Create(job).Error
}

// UpdateJobPosting updates an existing an existing job posting.
func (r *EmployerRepository) UpdateJobPosting(employerID, jobID string, job *models.Job) error {
	return r.db.Model(&models.Job{}).Where("id=? AND employer_id = ?", jobID, employerID).Updates(job).Error
}

// DeleteJobPosting deletes a job posting by ID.
func (r *EmployerRepository) DeleteJobPosting(employerID, jobID string) error {
	return r.db.Where("id = ? AND employer_id = ?", jobID, employerID).Delete(&models.Job{}).Error
}

// // CreateJobPosting creates a new job posting.
// func (r *EmployerRepository) CreateJobPosting(employerID string, job *models.Job) error {
// 	intEmployerID, _ := strconv.Atoi(employerID)
// 	uintEmployerID := uint(intEmployerID)
// 	job.EmployerID = uintEmployerID // Set the employer ID
// 	return r.db.Create(job).Error
// }

func (r *EmployerRepository) GetAllEmployers() ([]models.Employer, error) {
	var employers []models.Employer
	if err := r.db.Find(&employers).Error; err != nil {
		return nil, err
	}
	return employers, nil
}

func (r *EmployerRepository) GrpcListEmployers(ctx context.Context) ([]models.Employer, error) {
	var employers []models.Employer
	if err := r.db.WithContext(ctx).Find(&employers).Error; err != nil {
		return nil, err
	}
	return employers, nil
}

func (r *EmployerRepository) UpdateEmployer(employer *models.Employer) error {
	// Perform the update operation
	if err := r.db.Save(employer).Error; err != nil {
		return err
	}
	return nil
}
