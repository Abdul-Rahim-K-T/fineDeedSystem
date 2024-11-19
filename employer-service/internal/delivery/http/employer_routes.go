package http

import (
	"fineDeedSystem/employer-service/internal/middleware"
	"fineDeedSystem/employer-service/internal/usecase"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterEmployerRoutes(r *mux.Router, employerUsecase *usecase.EmployerUsecase) {
	handler := NewEmployerHandler(employerUsecase)

	r.HandleFunc("/employers/signup", handler.CreateEmployer).Methods("POST")
	r.HandleFunc("/employers/login", handler.Login).Methods("POST")
	r.Handle("/employers/logout", middleware.EmployerAuthMiddleware(employerUsecase)(http.HandlerFunc(handler.Logout))).Methods("GET")

	// Protected routes
	s := r.PathPrefix("employers").Subrouter()
	s.Use(middleware.EmployerAuthMiddleware(employerUsecase))
	s.HandleFunc("/myprofile", handler.GetEmployerProfile).Methods("GET")
	s.HandleFunc("/updateprofile", handler.UpdateEmployerProfile).Methods("PUT")
	s.HandleFunc("/changepassword", handler.ChangePassword).Methods("POST")
	s.HandleFunc("/dashboard", handler.GetEmployerDashboard).Methods("GET")

	// Job management
	s.HandleFunc("/jobs", handler.ListPostedJobs).Methods("GET")
	s.HandleFunc("/jobs", handler.CreateJobPosting).Methods("POST")

}
