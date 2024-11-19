package usecase

import (
	"encoding/json"
	"errors"
	"fineDeedSystem/applicant-service/internal/models"
	repository "fineDeedSystem/applicant-service/internal/repository"
	"fineDeedSystem/applicant-service/pkg/jwt"
	"fineDeedSystem/applicant-service/pkg/redisclient"
	"log"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type ApplicantUsecase struct {
	repo        repository.ApplicantRepository
	blacklist   map[string]bool
	redisClient *redisclient.RedisClient
}

func NewApplicantUsecase(repo repository.ApplicantRepository, redisClient *redisclient.RedisClient) *ApplicantUsecase {
	return &ApplicantUsecase{repo: repo, blacklist: make(map[string]bool), redisClient: redisClient}
}

func (u *ApplicantUsecase) CreateApplicant(applicant *models.Applicant) error {
	err := u.repo.Create(applicant)
	if err == nil {
		CacheApplicant(u.redisClient, applicant)
	}
	return err
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (u *ApplicantUsecase) Login(email, password string) (*models.Applicant, string, error) {
	applicant, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(applicant.Password), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(email, int(applicant.ID))
	if err != nil {
		return nil, "", err
	}
	return applicant, token, nil
}

// IsTokenValid checks if a given JWT token is valid
func (u *ApplicantUsecase) IsTokenValid(token string) bool {
	// Check if the token is blacklisted
	if _, blacklisted := u.blacklist[token]; blacklisted {
		return false
	}
	_, err := jwt.ValidateToken(token)
	return err == nil
}

func (u *ApplicantUsecase) Logout(token string) error {
	u.blacklist[token] = true
	return nil
}

func (u *ApplicantUsecase) GetProfile(token string) (*models.Applicant, error) {
	// Validate the token and extract the applicant ID (assuming your token contains the applicant ID)
	jwtToken, err := jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Extract UserID from token claims
	claims, ok := jwtToken.Claims.(jwtgo.MapClaims)
	if !ok || claims["UserID"] == nil {
		return nil, errors.New("invalid token claims")
	}

	// Convert the UserID from interface{} to int
	userID, ok := claims["UserID"].(float64) // or whatever the claim type is
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	// Convert to int since jwt-go uses float64 for numbers
	applicantID := int(userID)

	// Check Redis cache for applicant profile
	cacheKey := GetApplicantCacheKey(applicantID)
	applicantData, err := u.redisClient.Client.Get(u.redisClient.Client.Context(), cacheKey).Result()
	if err == nil && applicantData != "" {
		// If found in cache, return the profile
		// you might want unmarshal the cached data (if it's in JSON format) to a model
		var applicant models.Applicant
		if err := json.Unmarshal([]byte(applicantData), &applicant); err == nil {
			return &applicant, nil
		}

	}

	// Fetch the applicant's profile the repository
	applicant, err := u.repo.GetApplicantByID(applicantID)
	if err != nil {
		return nil, err
	}

	// Cache the profile in Redis with an expiration of 10 minutes
	err = CacheApplicant(u.redisClient, applicant)
	if err != nil {
		// Log if caching fails, but continue with response
		log.Printf("Error caching applicant profile: %v", err)
	}
	return applicant, nil
}

func (u *ApplicantUsecase) UpdateApplicant(applicant *models.Applicant) error {
	// Retrieve the existing applicant record from the database
	applicantID := int(applicant.ID)
	existingApplicant, err := u.repo.GetApplicantByID(applicantID)
	if err != nil {
		return errors.New("applicant not found")
	}

	// Update the applicant fields
	if applicant.FirstName != "" {
		existingApplicant.FirstName = applicant.FirstName
	}
	if applicant.LastName != "" {
		existingApplicant.LastName = applicant.LastName
	}
	if applicant.Email != "" {
		existingApplicant.Email = applicant.Email
	}
	if applicant.PhoneNumber != "" {
		existingApplicant.PhoneNumber = applicant.PhoneNumber
	}

	// Call the repository to update the applicant in the database
	if err := u.repo.UpdateApplicant(existingApplicant); err != nil {
		return errors.New("failed to update applicant")
	}

	// Update the cache updated profile in Redis
	err = CacheApplicant(u.redisClient, existingApplicant)
	if err != nil {
		//
		log.Printf("Error caching applicant profile: %v", err)
	}
	return nil
}

func (u *ApplicantUsecase) ChangePassword(token, oldPassword, newPassword string) error {
	jwtToken, err := jwt.ValidateToken(token)
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(jwtgo.MapClaims)
	if !ok || claims["UserID"] == nil {
		return errors.New("invalid token claims")
	}

	userID, ok := claims["UserID"].(float64)
	if !ok {
		return errors.New("invalid user ID in token")
	}
	applicantID := uint(userID)

	// Retrieves the existing applicant record from the database
	applicant, err := u.repo.GetApplicantByID(int(applicantID))
	if err != nil {
		return errors.New("applicant not found")
	}

	// Verify the old password
	if err := bcrypt.CompareHashAndPassword([]byte(applicant.Password), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old passoword")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update the password
	applicant.Password = string(hashedPassword)
	if err := u.repo.UpdateApplicant(applicant); err != nil {
		return errors.New("failed to update password")
	}

	// Update the cache with the new applicant data
	err = CacheApplicant(u.redisClient, applicant)
	if err != nil {
		log.Printf("Error caching applicant profile: %v", err)
	}
	return nil
}

func (u *ApplicantUsecase) DeleteApplicant(token string) error {
	// Validate the token and extract the applicant ID
	jwtToken, err := jwt.ValidateToken(token)
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(jwtgo.MapClaims)
	if !ok || claims["UserID"] == nil {
		return errors.New("invalid token claims")
	}

	userID, ok := claims["UserID"].(float64)
	if !ok {
		return errors.New("invalid usere ID in token")
	}

	applicantID := uint(userID)

	// Delete the applicant record from the database
	if err := u.repo.DeleteApplicant(applicantID); err != nil {
		return errors.New("failed to delete applicant")
	}

	// Remove the applicant data from the cache
	cacheKey := GetApplicantCacheKey(int(applicantID))
	err = u.redisClient.Client.Del(u.redisClient.Client.Context(), cacheKey).Err()
	if err != nil {
		log.Printf("Error removing applicant profile from cache: %v", err)
	}

	return nil
}
