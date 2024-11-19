package http

import (
	"context"
	"encoding/json"
	pb "fineDeedSystem/admin-service/proto/fineDeedSystem/proto/shared"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type GrpcHandler struct {
	GrpcClient  pb.EmployerServiceClient
	RedisClient *redis.Client
}

func NewGrpcHandler(grpcClient pb.EmployerServiceClient, redisClient *redis.Client) *GrpcHandler {
	return &GrpcHandler{GrpcClient: grpcClient, RedisClient: redisClient}
}

func (h *GrpcHandler) GrpcDeleteEmployer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employer ID", http.StatusBadRequest)
		return
	}

	req := &pb.DeleteEmployerRequest{
		Id: uint32(id),
	}

	// Call the gRPC client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.GrpcClient.DeleteEmployer(ctx, req)
	if err != nil {
		http.Error(w, "Failed to delete employer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate the cache for the deleted employer
	cacheKey := fmt.Sprintf("employer:%d", id)
	if err := h.RedisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		log.Printf("Failed to delete employer cache data: %v", err)
	}

	// Invalidate the cache for the list of all employers
	if err := h.RedisClient.Del(context.Background(), "employers").Err(); err != nil {
		log.Printf("Failed to delete employers list cache: %v", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (h *GrpcHandler) GrpcUpdateEmployer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employer ID", http.StatusBadRequest)
		return
	}

	var employerRequest pb.Employer
	if err := json.NewDecoder(r.Body).Decode(&employerRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the ID from the URL parameter
	employerRequest.Id = uint32(id)

	// Check Redis cache first and invalidate if exists
	cacheKey := fmt.Sprintf("employer:%d", id)
	log.Printf("Attempting to retrieve employer data from Redis with key: %s", cacheKey)

	// Check if employer data is cached
	_, err = h.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.Printf("Employer data found in Redis, invalidating cache for employer ID: %d", id)
		// Invalidate the cache by deleting the cached data
		err := h.RedisClient.Del(context.Background(), cacheKey).Err()
		if err != nil {
			log.Printf("Error invalidating Redis cache: %v", err)
		}
	} else if err != redis.Nil {
		log.Printf("Error checking Redis cache: %v", err)
	}

	// Call the gRPC method to update the employer
	req := &pb.UpdateEmployerRequest{Employer: &employerRequest}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.GrpcClient.UpdateEmployer(ctx, req)
	if err != nil {
		http.Error(w, "Failed to update employer", http.StatusInternalServerError)
		return
	}

	// Update the Redis cache with the updated employer data
	updatedEmployerJSON, err := json.Marshal(res.Employer)
	if err != nil {
		log.Printf("Error marshalling updated employer to JSON: %v", err)
		http.Error(w, "Failed to serialize updated employer", http.StatusInternalServerError)
		return
	}

	// set the updated employer data in Redis (with a new expiration time)
	err = h.RedisClient.Set(context.Background(), cacheKey, updatedEmployerJSON, time.Minute*10).Err()
	if err != nil {
		log.Printf("Failed to cache updated employer data: %v", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res.Employer)
}

func (h *GrpcHandler) GrpcGetEmployerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employer ID", http.StatusBadRequest)
		return
	}

	// Check Redis cache first
	cacheKey := fmt.Sprintf("employer:%d", employerID)
	log.Printf("Attempting to retrieve employer data from Redist with key: %s", cacheKey)
	cachedEmployer, err := h.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.Printf("Employer data retrieved from Redis with key: %s", cachedEmployer)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedEmployer))
		return
	} else if err != redis.Nil {
		log.Printf("Error retrieving employer data from Redis: %v", err)
	}

	req := &pb.GetEmployerByIDRequest{Id: uint32(employerID)}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.GrpcClient.GetEmployerByID(ctx, req)
	if err != nil {
		http.Error(w, "Failed to get employer", http.StatusInternalServerError)
		return
	}

	employerJSON, err := json.Marshal(res.Employer)
	if err != nil {
		http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
		return
	}

	// Cache the response in Redis
	err = h.RedisClient.Set(context.Background(), cacheKey, employerJSON, time.Minute*10).Err()
	if err != nil {
		log.Printf("Failed to cache employer data: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(employerJSON)
}

func (g *GrpcHandler) GrpcGetAllEmployers(w http.ResponseWriter, r *http.Request) {
	log.Println("GrpcGetAllEmployers: Request recieved")

	// Check if data is available in Redis
	cachedEmployers, err := g.RedisClient.Get(context.Background(), "employers").Result()
	if err == nil {
		log.Println("Data retrieved from Redis cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedEmployers))
		return
	}

	// Make gRPC call to get the list of employers
	req := &pb.ListEmployersRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := g.GrpcClient.ListEmployers(ctx, req)
	if err != nil {
		log.Printf("Error retrieving employers via gRPC: %v", err)
		http.Error(w, "Failed to retrieve employers via gRPC", http.StatusInternalServerError)
		return
	}

	// var employers []models.Employer
	// for _, e := range res.Employers {
	// 	employers = append(employers, models.Employer{
	// 		ID:          uint(e.Id),
	// 		Name:        e.Name,
	// 		Email:       e.Email,
	// 		Phone:       e.Phone,
	// 		CompanyName: e.CompanyName,
	// 	})
	// }
	// log.Printf("Employers retrieved via gRPC: %+v", employers)

	// Serialize the response to JSON
	employersJSON, err := json.Marshal(res.Employers)
	if err != nil {
		log.Printf("Error marshaling employers to JSON: %v", err)
		http.Error(w, "Failed to serialize employers", http.StatusInternalServerError)
		return
	}

	// Cache the data in Redis
	// jsonEmployers, err := json.Marshal(employers)
	err = g.RedisClient.Set(context.Background(), "employers", employersJSON, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Error caching data in Redis: %v", err)
	} else {
		log.Println("Data cached in Redis")
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(employersJSON)
}

func (h *AdminHandler) GrpcCreateEmployers(w http.ResponseWriter, r *http.Request) {
	log.Println("GrpcCreateEmployers invoked ?grpc_handler.go file?")
	var employerRequest pb.Employer
	if err := json.NewDecoder(r.Body).Decode(&employerRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Log individual fields to avoid copying the mutex
	log.Printf("Decoded employerRequest: Name=%s, Email=%s, Phone=%s, CompanyName=%s, Password=%s",
		employerRequest.Name, employerRequest.Email, employerRequest.Phone, employerRequest.CompanyName, employerRequest.Password)

	// Check if the user already exists in Redis
	existingEmployer, err := h.RedisClient.Get(context.Background(), "employer_"+employerRequest.Email).Result()
	if err == nil && existingEmployer != "" {
		log.Println("Error: Employer already exists in Redis")
		http.Error(w, "Employer already exists in database", http.StatusConflict)
		return
	}

	// Call the gRPC method to create the employer
	response, err := h.usecase.GrpcCreateEmployer(r.Context(), &employerRequest)
	if err != nil {
		log.Printf("Error creating employer: %v", err)
		http.Error(w, "Failed to create employer", http.StatusInternalServerError)
		return
	}

	// Invalidate the cache in Redis(because the actual data are updated from above so we need to delete all the cached data to ensure updated data shou
	// be get upcoming retreiving of data from the cache data)
	if err := h.RedisClient.Del(context.Background(), "employers").Err(); err != nil {
		log.Printf("Error invalidating Redis cache: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
