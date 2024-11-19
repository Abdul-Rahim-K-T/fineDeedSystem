package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fineDeedSystem/employer-service/internal/models"
	"fineDeedSystem/employer-service/internal/repository/postgres"
	"fineDeedSystem/employer-service/pkg/database"
	"fineDeedSystem/employer-service/pkg/jwt"
	"fineDeedSystem/employer-service/proto/fineDeedSystem/proto/shared"
	"fmt"

	"log"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

type EmployerUsecase struct {
	repo        *postgres.EmployerRepository
	redisClient *database.RedisClient
	queueName   string
	shared.UnimplementedEmployerServiceServer
}

func NewEmployerUsecase(repo *postgres.EmployerRepository, redisClient *database.RedisClient, queueName string) *EmployerUsecase {
	return &EmployerUsecase{repo: repo, redisClient: redisClient, queueName: queueName}
}

func (u *EmployerUsecase) DeleteEmployer(ctx context.Context, id uint32) error {
	// Delete the employer from the repository
	err := u.repo.DeleteEmployer(uint(id))
	if err != nil {
		return err
	}

	// Invalidate the cache for the deleted employer
	cacheKey := fmt.Sprintf("employer:%d", id)
	if err := u.redisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		log.Printf("Failed to delete employer cache data: %v", err)
	}

	// Invalidate the cache for the list of all employers
	if err := u.redisClient.Del(context.Background(), "employers").Err(); err != nil {
		log.Printf("Failed to delete employers list cache: %v", err)
	}
	return nil
}

func (u *EmployerUsecase) UpdateEmployer(ctx context.Context, employer *models.Employer) (*models.Employer, error) {
	// // Validate input
	// if employer.ID == 0 || employer.Name == "" || employer.Email == "" || employer.Phone == "" || employer.CompanyName == "" {
	// 	return nil, errors.New("all fields are required")
	// }

	// check if an employer with the give ID exists
	existingEmployer, err := u.repo.GetEmployerByID(uint(employer.ID))
	if err != nil {
		return nil, err
	}
	if existingEmployer == nil {
		return nil, errors.New("employer not found")
	}

	// Initialize a flag to track if any field is changed
	fieldsChanged := false

	// Update the fields only if they are provided and different from the existing ones
	if employer.Name != "" && employer.Name != existingEmployer.Name {
		fmt.Println("Name changed from", existingEmployer.Name, "to", employer.Name)
		existingEmployer.Name = employer.Name
		fieldsChanged = true
	}
	if employer.Email != "" && employer.Email != existingEmployer.Email {
		fmt.Println("Email changed from", existingEmployer.Email, "to", employer.Email)
		existingEmployer.Email = employer.Email
		fieldsChanged = true
	}
	if employer.Phone != "" && employer.Phone != existingEmployer.Phone {
		fmt.Println("Phone changed from ", existingEmployer.Phone, "to", employer.Phone)
		existingEmployer.Phone = employer.Phone
		fieldsChanged = true
	}
	if employer.CompanyName != "" && employer.CompanyName != existingEmployer.CompanyName {
		fmt.Println("CompanyName changed from", existingEmployer.CompanyName, "to", employer.CompanyName)
		existingEmployer.CompanyName = employer.CompanyName
		fieldsChanged = true
	}
	if employer.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employer.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		// Only set fieldsChanged if the hashed password is different
		if existingEmployer.Password != string(hashedPassword) {
			fmt.Println("Password changed")
			existingEmployer.Password = string(hashedPassword)
			fieldsChanged = true
		}
	}

	// If no fields are changed, return an error
	if !fieldsChanged {
		fmt.Println("No fields changed")
		return nil, errors.New("no fields are changed")
	}

	// // Hash the password if it's provided
	// if employer.Password != "" {
	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employer.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	employer.Password = string(hashedPassword)
	// } else {
	// 	employer.Password = existingEmployer.Password
	// }

	// Update the employer in the repository
	if err := u.repo.UpdateEmployer(existingEmployer); err != nil {
		return nil, err
	}

	return existingEmployer, nil
}

func (u *EmployerUsecase) GetEmployerByID(ctx context.Context, id uint32) (*models.Employer, error) {
	return u.repo.GetEmployerByID(uint(id))
}

// CreateEmployer handles the logic for creating a new employer
func (u *EmployerUsecase) CreateEmployerLogic(employer *models.Employer) (*models.Employer, error) {

	log.Printf("Received employer details: Name=%s, Email=%s, Phone=%s, CompanyName=%s, Password=%s",
		employer.Name, employer.Email, employer.Phone, employer.CompanyName, employer.Password)
	// Validate input
	if employer.Name == "" || employer.Email == "" || employer.Phone == "" || employer.CompanyName == "" || employer.Password == "" {
		return nil, errors.New("all fields are required")
	}

	// Check if an employer with same already exists
	existingEmployer, err := u.repo.FindByEmail(employer.Email)
	if err != nil {
		return nil, err
	}
	if existingEmployer != nil {
		return nil, errors.New("an employer with this email already exists")
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employer.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	employer.Password = string(hashedPassword)

	// Call the repository to save the employer
	if err := u.repo.CreateEmployer(employer); err != nil {
		return nil, err
	}

	// Cache the newly created employer data
	cacheKey := fmt.Sprintf("employer:%d", employer.ID)
	employerJSON, err := json.Marshal(employer)
	if err != nil {
		log.Printf("Failed to marshal employer data: %v", err)
	} else {
		if err := u.redisClient.Set(cacheKey, employerJSON, 0).Err(); err != nil {
			log.Printf("Failed to cache employer data: %v", err)
		}
	}

	return employer, nil
}

func (u *EmployerUsecase) Login(creds jwt.Credentials) (string, error) {
	fmt.Println("Attempting login for:", creds.Username)

	employer, err := u.repo.GetByEmail(creds.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(employer.Password), []byte(creds.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(employer.Email, true) // true indicates employer role
	if err != nil {
		return "", err
	}

	// Define the token expiration time
	expiration := 24 * time.Hour

	// Store the token in Redis
	err = u.StoreToken(token, expiration)
	if err != nil {
		log.Println("Store token in Redis:", err)
		return "", err
	}

	return token, nil
}

func (u *EmployerUsecase) Logout(token string) error {
	err := u.redisClient.Set(token, "blacklisted", 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *EmployerUsecase) GetEmployerProfile(employerID string) (*models.Employer, error) {
	return u.repo.GetEmployerProfile(employerID)
}

func (u *EmployerUsecase) UpdateEmployerProfile(employerID string, updatedProfile *models.Employer) error {
	return u.repo.UpdateEmployerProfile(employerID, updatedProfile)
}

func (u *EmployerUsecase) ChangePassword(employerID, oldPassword, newPassword string) error {
	return u.repo.ChangePassword(employerID, oldPassword, newPassword)
}

func (u *EmployerUsecase) GetEmployerDashboard(employerID string) (interface{}, error) {
	return u.repo.GetEmployerDashboard(employerID)
}

func (u *EmployerUsecase) ListPostedJobs(employerID string) ([]models.Job, error) {
	return u.repo.ListPostedJobs(employerID)
}

func (u *EmployerUsecase) CreateJobPosting(employerID string, job *models.Job) error {
	return u.repo.CreateJobPosting(employerID, job)
}

func (u *EmployerUsecase) UpdateJobPosting(employerID, jobID string, job *models.Job) error {
	return u.repo.UpdateJobPosting(employerID, jobID, job)
}

func (u *EmployerUsecase) DeleteJobPosting(employerID, jobID string) error {
	return u.repo.DeleteJobPosting(employerID, jobID)
}

func (u *EmployerUsecase) GetAllEmployers() ([]models.Employer, error) {
	return u.repo.GetAllEmployers()
}

func (u *EmployerUsecase) IsTokenBlacklisted(token string) (bool, error) {
	val, err := u.redisClient.Get(token).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Token not found in blacklist
		}
		return false, err
	}
	return val == "blacklisted", nil
}

// func (u *EmployerUsecase) ListEmployers(ctx context.Context, req *shared.ListEmployersRequest) (*shared.ListEmployersResponse, error) {
// 	employers, err := u.repo.GrpcListEmployers(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var employerList []*shared.Employer
// 	for _, emp := range employers {
// 		employerList = append(employerList, &shared.Employer{
// 			Id:          uint32(emp.ID),
// 			Name:        emp.Name,
// 			Email:       emp.Email,
// 			Phone:       emp.Phone,
// 			CompanyName: emp.CompanyName,
// 		})
// 	}
// 	return &shared.ListEmployersResponse{Employers: employerList}, nil
// }

func (u *EmployerUsecase) StoreToken(token string, expiration time.Duration) error {
	err := u.redisClient.Set(token, "active", expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
