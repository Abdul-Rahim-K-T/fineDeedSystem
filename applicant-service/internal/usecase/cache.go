package usecase

import (
	"encoding/json"
	"fineDeedSystem/applicant-service/internal/models"
	"fineDeedSystem/applicant-service/pkg/redisclient"
	"log"
	"strconv"
	"time"
)

// Generates the cache key for an applicant profile
func GetApplicantCacheKey(applicantID int) string {
	return "applicant_profile_" + strconv.Itoa(applicantID)
}

// Caches the applicant profile in Redis
func CacheApplicant(redisClient *redisclient.RedisClient, applicant *models.Applicant) error {
	applicantID := int(applicant.ID)
	cacheKey := GetApplicantCacheKey(applicantID)

	applicantJSON, err := json.Marshal(applicant)
	if err != nil {
		return err
	}

	err = redisClient.Client.Set(redisClient.Client.Context(), cacheKey, applicantJSON, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Error caching applicant profile: %v", err)
		return err
	}
	return nil
}

func GetCachedApplicant(redisClient *redisclient.RedisClient, applicantID int) (*models.Applicant, error) {
	cacheKey := GetApplicantCacheKey(applicantID)

	applicantData, err := redisClient.Client.Get(redisClient.Client.Context(), cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var applicant models.Applicant
	err = json.Unmarshal([]byte(applicantData), &applicant)
	if err != nil {
		return nil, err
	}

	return &applicant, nil
}
