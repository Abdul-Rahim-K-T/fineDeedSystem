package http

import (
	"encoding/json"
	"fineDeedSystem/applicant-service/internal/models"
	"fineDeedSystem/applicant-service/internal/usecase"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ApplicantHandler struct {
	Usecase *usecase.ApplicantUsecase
}

func NewApplicantHandler(uc *usecase.ApplicantUsecase) *ApplicantHandler {
	return &ApplicantHandler{Usecase: uc}
}

func (h *ApplicantHandler) CreateApplicant(w http.ResponseWriter, r *http.Request) {
	// // Implement the logic for creating an applicant
	// var applicant models.Applicant
	// if err := json.NewDecoder(r.Body).Decode(&applicant); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Parse form-data
	err := r.ParseMultipartForm(10 << 20) // Limit file size of to 10 MB
	if err != nil {
		http.Error(w, "Unble to parse form data", http.StatusBadRequest)
		log.Println("ParseMultipartForm error:", err)
		return
	}

	// Get other fields
	applicant := models.Applicant{
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Email:       r.FormValue("email"),
		PhoneNumber: r.FormValue("phone_number"),
	}

	// Hash the password
	password := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		log.Println("GenerateFromPassword error:", err)
		return
	}
	applicant.Password = string(hashedPassword)

	// Handle profile_picture upload
	profilePicFile, profilePicHeader, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "Profile picture is required", http.StatusBadRequest)
		log.Println("FormFile profile_picture error:", err)
		return
	}
	defer profilePicFile.Close()

	// Ensure the directory exists
	profilePicDir := filepath.Join("uploads", "profile_pictures")
	if err := os.MkdirAll(profilePicDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create directory for profile pictures", http.StatusInternalServerError)
		log.Println("MkdirAll profile_picture directory error:", err)
		return
	}

	profilePicPath := filepath.Join(profilePicDir, profilePicHeader.Filename)
	profilePicDst, err := os.Create(profilePicPath)
	if err != nil {
		http.Error(w, "Unable to save profile picture", http.StatusInternalServerError)
		log.Println("Create profile_picture file error:", err)
		return
	}
	defer profilePicDst.Close()
	_, err = io.Copy(profilePicDst, profilePicFile)
	if err != nil {
		http.Error(w, "Unable to save profile picture", http.StatusInternalServerError)
		log.Println("Copy profile_picture file error:", err)
		return
	}
	applicant.ProfilePicture = profilePicPath

	// Handle resume upload
	resumeFile, resumeHeader, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Resume is required", http.StatusBadRequest)
		log.Println("FormFile resume error:", err)
		return
	}
	defer resumeFile.Close()

	// Ensure the directory exists
	resumeDir := filepath.Join("uploads", "resumes")
	if err := os.MkdirAll(resumeDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create directory for resumes", http.StatusInternalServerError)
		log.Println("MkdirAll resume directive error:", err)
		return
	}

	resumePath := filepath.Join(resumeDir, resumeHeader.Filename)
	resumeDst, err := os.Create(resumePath)
	if err != nil {
		http.Error(w, "Unable to save resume", http.StatusInternalServerError)
		log.Println("Create resume file error:", err)
		return
	}
	defer resumeDst.Close()
	_, err = io.Copy(resumeDst, resumeFile)
	if err != nil {
		http.Error(w, "Unable to save resume", http.StatusInternalServerError)
		log.Println("Copy resume file error:", err)
		return
	}
	applicant.Resume = resumePath

	// Save applicant using the usecase
	if err := h.Usecase.CreateApplicant(&applicant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("CreateApplicant usecase error:", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Applicant created"})
}

type ApplicantResponse struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	Resume         string `json:"resume"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func (h *ApplicantHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Decode error:", err)
		return
	}

	// Use the usecase to authenticate the user
	applicant, token, err := h.Usecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// set the token in the response header
	w.Header().Set("Authorization", "Bearer "+token)

	// Set the token in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// applicant response
	applicantResp := ApplicantResponse{
		ID:             int(applicant.ID),
		FirstName:      applicant.FirstName,
		LastName:       applicant.LastName,
		Email:          applicant.Email,
		PhoneNumber:    applicant.PhoneNumber,
		ProfilePicture: applicant.ProfilePicture,
		Resume:         applicant.Resume,
		CreatedAt:      applicant.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:      applicant.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	// Include applicant details in the response
	response := map[string]interface{}{
		"token":      token,
		"applicants": applicantResp,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ApplicantHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the token from the cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Ensuere this is true if we are using HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Expired token to remove it
	})

	// You can also invalidate the JWT on the server if necessary  (blacklist)
	token := r.Header.Get("Authorization")
	if token == "" {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token = cookie.Value
		} else {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}
	}

	token = strings.TrimPrefix(token, "Bearer ")

	if err := h.Usecase.Logout(token); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear the token from the cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

func (h *ApplicantHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the header or cookie
	token := r.Header.Get("Authorization")
	if token == "" {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token = cookie.Value
		} else {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}
	}

	token = strings.TrimPrefix(token, "Bearer ")

	// Use the usecase to get the applicant profile
	applicant, err := h.Usecase.GetProfile(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response struct
	applicantResp := ApplicantResponse{
		ID:             int(applicant.ID),
		FirstName:      applicant.FirstName,
		LastName:       applicant.LastName,
		Email:          applicant.Email,
		PhoneNumber:    applicant.PhoneNumber,
		ProfilePicture: applicant.ProfilePicture,
		Resume:         applicant.Resume,
		CreatedAt:      applicant.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:      applicant.UpdatedAt.Format("2006-01-20T15:04:05.000Z"),
	}

	// Return the profile in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(applicantResp)
}

func (h *ApplicantHandler) UpdateApplicant(w http.ResponseWriter, r *http.Request) {
	var applicant models.Applicant
	if err := json.NewDecoder(r.Body).Decode(&applicant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Usecase.UpdateApplicant(&applicant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (h *ApplicantHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var changepasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&changepasswordRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")
	if err := h.Usecase.ChangePassword(token, changepasswordRequest.OldPassword, changepasswordRequest.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ApplicantHandler) DeleteApplicant(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if err := h.Usecase.DeleteApplicant(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
