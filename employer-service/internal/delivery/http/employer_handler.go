package http

import (
	"encoding/json"
	"fineDeedSystem/employer-service/internal/models"
	"fineDeedSystem/employer-service/internal/usecase"
	"fineDeedSystem/employer-service/pkg/constants"
	"fineDeedSystem/employer-service/pkg/jwt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// EmployerHandler handles HTTP requests related to employers.
type EmployerHandler struct {
	Usecase *usecase.EmployerUsecase
}

// NewEmployerHandler creates a new EmployerHandler
func NewEmployerHandler(usecase *usecase.EmployerUsecase) *EmployerHandler {
	return &EmployerHandler{Usecase: usecase}
}

// CreateEmployer handles the creation of a new employer.
func (h *EmployerHandler) CreateEmployer(w http.ResponseWriter, r *http.Request) {
	var employer models.Employer

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&employer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the use case to create the employer
	createdEmployer, err := h.Usecase.CreateEmployerLogic(&employer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respod with the created employer
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEmployer)
}

// // GetEmployByID handles the retrieval of an employer by ID.
// func (h *EmployerHandler) GetEmployByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		http.Error(w, "Missing employer ID", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid employer ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Call the use case to get the employer by ID
// 	employer, err := h.Usecase.GetEmployerByID(id)
// 	if err != nil {
// 		if err == h.Usecase.ErrEmployerNotFound {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Respond with the employer data
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(employer)
// }

// Login handles employer login.
func (h *EmployerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds jwt.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.Usecase.Login(creds)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // This helps mitigates XSS attacks
		Secure:   true, // This ensures the cookie is only sent through HTTPS
		Path:     "/",
	})

	// Optionally set the Authorization header
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Logout handles employer logout.
func (h *EmployerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(constants.UserClaimsKey).(*jwt.Claims)
	if !ok || claims == nil {
		// log.Printf()
		http.Error(w, "No claims found", http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	var tokenString string
	if authHeader != "" && len(authHeader) > len("Bearer ") {
		tokenString = authHeader[len("Bearer "):]
	} else {
		// If token is not in Authorization header, check the cookie
		cookie, err := r.Cookie("jwt_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Authorization header is missing or invalid", http.StatusUnauthorized)
			return
		}
		tokenString = cookie.Value
	}

	if err := h.Usecase.Logout(tokenString); err != nil {
		http.Error(w, "Error blacklisting token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func (h *EmployerHandler) GetEmployerProfile(w http.ResponseWriter, r *http.Request) {
	// Assuming you have a way to get the employer ID from the JWT token
	employerID := r.Context().Value("employer_id").(string)
	profile, err := h.Usecase.GetEmployerProfile(employerID)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *EmployerHandler) UpdateEmployerProfile(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	var updatedProfile models.Employer
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Usecase.UpdateEmployerProfile(employerID, &updatedProfile)
	if err != nil {
		http.Error(w, "Unble to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *EmployerHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	var passwordChangeRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&passwordChangeRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Usecase.ChangePassword(employerID, passwordChangeRequest.OldPassword, passwordChangeRequest.NewPassword)
	if err != nil {
		http.Error(w, "Unble to change password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EmployerHandler) GetEmployerDashboard(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	dashboard, err := h.Usecase.GetEmployerDashboard(employerID)
	if err != nil {
		http.Error(w, "Dashboard data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *EmployerHandler) ListPostedJobs(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	jobs, err := h.Usecase.ListPostedJobs(employerID)
	if err != nil {
		http.Error(w, "Unable to fetch jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *EmployerHandler) CreateJobPosting(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Usecase.CreateJobPosting(employerID, &job)
	if err != nil {
		http.Error(w, "Unble to create job", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *EmployerHandler) UpdateJobPosting(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	jobID := mux.Vars(r)["id"]
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Usecase.UpdateJobPosting(employerID, jobID, &job)
	if err != nil {
		http.Error(w, "Unable to update job", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *EmployerHandler) DeleteJobPosting(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value("employer_id").(string)
	jobID := mux.Vars(r)["id"]

	err := h.Usecase.DeleteJobPosting(employerID, jobID)
	if err != nil {
		http.Error(w, "Unable to delete job ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
