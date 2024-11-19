package http

import (
	"fineDeedSystem/admin-service/internal/models"
	"fineDeedSystem/admin-service/internal/usecase"
	"fineDeedSystem/admin-service/pkg/jwt"
	"fmt"

	"time"

	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type AdminHandler struct {
	usecase     *usecase.AdminUsecase
	RedisClient *redis.Client
}

func NewAdminHandler(usecase *usecase.AdminUsecase, redisClient *redis.Client) *AdminHandler {
	return &AdminHandler{
		usecase:     usecase,
		RedisClient: redisClient,
	}
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin

	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.usecase.CreateAdmin(admin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) GetAdminByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	admin, err := h.usecase.FindAdminByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Adminname string `json:"adminname"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	token, err := h.usecase.Login(credentials.Adminname, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Prevents Javascript access to the cookie
		MaxAge:   86400,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully"))
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Check for token in the Cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Error retrieving cookie:", err)
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	fmt.Println("Retrieved token from cookie:", token)

	// Validate the token
	claims, err := jwt.ValidateToken(token)
	if err != nil {
		fmt.Println("Error validating token:", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Blacklist the token
	expiration := time.Until(time.Unix(claims.ExpiresAt, 0))
	if err := h.usecase.BlacklistToken(token, expiration); err != nil {
		fmt.Println("Error blacklisting token:", err)
		http.Error(w, "Failed to blacklist token", http.StatusUnauthorized)
		return
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete the cookie
		HttpOnly: true,
	})

	fmt.Println("Successfully logged out")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

// func (h *AdminHandler) GrpcGetAllEmployers(w http.ResponseWriter, r *http.Request) {
// 	log.Println("GrpcGetAllEmployers: Request recieved")

// 	// Call the gRPC ListEmployers method
// 	req:= &pb.ListEmployersRequest{}
// 	res, err:= h.GrpcClient.ListEmployers(context.Background(),req)
// 	if err != nil {
// 		log.Printf("Error calling gRPC ListEmployers: %v", err)
// 		http.Error(w, "Failed to list employers via gRPC", http.StatusInternalServerError)
// 		return
// 	}

// 	// Convert the gRPC response to the format used in your project
// 	var employers []models.Employer
// 	for _, e:= range res.Employers {
// 		employers = append(employers, models.Employer{
// 			ID: uint(e.Id),
// 			Name: e.Name,
// 			Email: e.Email,
// 			Phone: e.Phone,
// 			CompanyName: e.CompanyName,
// 		})
// 	}
// 	log.Printf("Employers retrieved via gRPC: %+v", employers)
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(employers)

// }

func (h *AdminHandler) GetAllEmployers(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	// Log the entire URL and query parameters
	query := r.URL.Query()
	fmt.Println("Query parameters:", query)
	employers, err := h.usecase.GetAllEmployers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employers)
}

// func (h *AdminHandler) GrpcCreateEmployers(w http.ResponseWriter, r *http.Request) {
// 	var employerRequest pb.Employer
// 	if err := json.NewDecoder(r.Body).Decode(&employerRequest); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	// Call the gRPC method to create the employer
// 	response, err := h.usecase.GrpcCreateEmployer(r.Context(), &employerRequest)
// 	if err != nil {
// 		http.Error(w, "Failed to create employer", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }
