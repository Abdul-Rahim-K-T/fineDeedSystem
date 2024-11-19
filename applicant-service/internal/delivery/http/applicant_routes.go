package http

import (
	"fineDeedSystem/applicant-service/internal/middleware"
	"fineDeedSystem/applicant-service/internal/usecase"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterApplicantRoutes(r *mux.Router, applicantUsecase *usecase.ApplicantUsecase) {
	handler := NewApplicantHandler(applicantUsecase)

	r.HandleFunc("/applicants/signup", handler.CreateApplicant).Methods("POST")
	r.HandleFunc("/applicants/login", handler.Login).Methods("POST")
	r.Handle("/applicants/logout", middleware.ApplicantAuthMiddleware(applicantUsecase)(http.HandlerFunc(handler.Logout))).Methods("GET")

	r.Handle("/applicants/profile", middleware.ApplicantAuthMiddleware(applicantUsecase)(http.HandlerFunc(handler.GetProfile))).Methods("GET")

	r.Handle("/applicants/profile", middleware.ApplicantAuthMiddleware(applicantUsecase)(http.HandlerFunc(handler.UpdateApplicant))).Methods("PUT")
	r.Handle("/applicants/change-password", middleware.ApplicantAuthMiddleware(applicantUsecase)(http.HandlerFunc(handler.ChangePassword))).Methods("POST")
	r.Handle("/applicants/delete", middleware.ApplicantAuthMiddleware(applicantUsecase)(http.HandlerFunc(handler.DeleteApplicant))).Methods("DELETE")

}
