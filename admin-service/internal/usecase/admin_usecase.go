package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"fineDeedSystem/admin-service/internal/models"
	"fineDeedSystem/admin-service/internal/repository/postgres"
	"fineDeedSystem/admin-service/pkg/jwt"

	"time"

	"fineDeedSystem/admin-service/proto/fineDeedSystem/proto/shared"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	repo            postgres.AdminRepository
	redisClient     *redis.Client
	rabbitMQUsecase *EmployerRabbitMQUsecase
	grpcUsecase     *EmployerGrpcUsecase
}

func NewAdminUsecase(repo postgres.AdminRepository, redisClient *redis.Client, rabbitMQUsecase *EmployerRabbitMQUsecase, grpcUsecase *EmployerGrpcUsecase) *AdminUsecase {
	return &AdminUsecase{
		repo:            repo,
		redisClient:     redisClient,
		rabbitMQUsecase: rabbitMQUsecase,
		grpcUsecase:     grpcUsecase,
	}
}

func (a *AdminUsecase) InvalidateCache(message []byte) {
	// Assuming the message contains an employer ID or some identifier
	// Parse the message as needed
	employerID := string(message)

	// Remove the cached data
	cacheKey := "employer:" + employerID
	err := a.redisClient.Del(context.Background(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache key %s: %v", cacheKey, err)
	} else {
		log.Printf("Cache invalidated for key %s", cacheKey)
	}
}

func (u *AdminUsecase) CreateAdmin(admin models.Admin) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)
	return u.repo.CreateAdmin(admin)
}

func (u *AdminUsecase) FindAdminByID(id uint) (models.Admin, error) {
	return u.repo.FindAdminByID(id)
}

func (u *AdminUsecase) Login(adminname, password string) (string, error) {
	admin, err := u.repo.FindAdminByAdminname(adminname)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := jwt.GenerateToken(admin.Adminname)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *AdminUsecase) BlacklistToken(token string, exp time.Duration) error {
	return u.redisClient.Set(context.Background(), token, true, exp).Err()
}

func (u *AdminUsecase) IsTokenBlacklisted(token string) (bool, error) {
	result, err := u.redisClient.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "true", nil
}

// func (u *AdminUsecase) GetAllEmployers() ([]*employerpb.Employer, error) {
// 	ctx := context.Background()
// 	log.Println("Calling gRPC client GetAllEmployers")
// 	employers, err := u.grpcClient.GetAllEmployers(ctx, &employerpb.Empty{})
// 	if err != nil {
// 		log.Printf("Error from gRPC client: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("Error from gRPC client: %+v", employers)
// 	return employers.Employer, nil
// }

func (u *AdminUsecase) GetAllEmployers() ([]*models.Employer, error) {

	ctx := context.Background()
	cacheKey := "employers_cache"

	// Try to fetch from cache
	cachedEmployers, err := u.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var employers []*models.Employer
		if err := json.Unmarshal([]byte(cachedEmployers), &employers); err == nil {
			log.Println("Cache hit for employers")
			return employers, nil
		}
	}

	// If cache miss or error, fetch from RabbitMQ
	employers, err := u.rabbitMQUsecase.GetAllEmployers()
	if err != nil {
		return nil, err
	}

	// Cache the result
	employersJSON, err := json.Marshal(employers)
	if err == nil {
		u.redisClient.Set(ctx, cacheKey, employersJSON, 5*time.Minute)
	}
	return employers, nil
}

func (u *AdminUsecase) GrpcCreateEmployer(ctx context.Context, employer *shared.Employer) (*shared.Employer, error) {

	log.Println("GrpcCreateEmployer at admin_usecase.go file ")
	return u.grpcUsecase.GrpcCreateEmployer(ctx, employer)
}
