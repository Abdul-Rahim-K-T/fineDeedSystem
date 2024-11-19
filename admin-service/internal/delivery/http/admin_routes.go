package http

import (
	"fineDeedSystem/admin-service/internal/usecase"
	"fmt"
	"os"

	"fineDeedSystem/admin-service/internal/middleware"

	pb "fineDeedSystem/admin-service/proto/fineDeedSystem/proto/shared"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router, adminUsecase *usecase.AdminUsecase, grpcClient pb.EmployerServiceClient, redisClient *redis.Client) {
	handler := NewAdminHandler(adminUsecase, redisClient)
	grpcHandler := NewGrpcHandler(grpcClient, redisClient)

	r.HandleFunc("/admin", handler.CreateAdmin).Methods("POST")
	r.HandleFunc("/admin/{id:[0-9]+}", handler.GetAdminByID).Methods("GET")
	// Admin Authentication Routes
	r.HandleFunc("/admin/login", handler.Login).Methods("post")
	r.HandleFunc("/admin/logout", middleware.AuthMiddleware(handler.Logout)).Methods("POST")

	// Employer Management Routes
	// r.HandleFunc("/admin/employers", middleware.AuthMiddleware(handler.GetAllEmployers)).Methods("GET")
	r.HandleFunc("/admin/ListEmployers", middleware.AuthMiddleware(handler.GetAllEmployers)).Methods("GET")
	r.HandleFunc("/admin/GrpcListEmployers", middleware.AuthMiddleware(grpcHandler.GrpcGetAllEmployers)).Methods("GET")

	// r.HandleFunc("/admin/employers", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Received request to get all employers")
	// 	handler.GetAllEmployers(w, r)
	// })).Methods("GET")

	if os.Getenv("ENABLE_EMPLOYER_CREATION") == "true" {
		fmt.Println("Registering /admin/CreateEmployer route")
		r.HandleFunc("/admin/CreateEmployer", middleware.AuthMiddleware(handler.GrpcCreateEmployers)).Methods("POST")
	}
	r.HandleFunc("/admin/employers/{id:[0-9]+}", middleware.AuthMiddleware(grpcHandler.GrpcGetEmployerByID)).Methods("GET")
	r.HandleFunc("/admin/employers/{id:[0-9]+}", middleware.AuthMiddleware(grpcHandler.GrpcUpdateEmployer)).Methods("PUT")
	r.HandleFunc("/admin/employers/{id:[0-9]+}", middleware.AuthMiddleware(grpcHandler.GrpcDeleteEmployer)).Methods("DELETE")

	// // Applicant Management Routes
	// r.HandleFunc("/admin/applicants", middleware.AuthMiddleware(handler.GetAllApplicants)).Methods("GET")
	// if os.Getenv("ENABLE_APPLICANT_CREATION") == "true" {
	// 	r.HandleFunc("/admin/applicants", middleware.AuthMiddleware(handler.CreateApplicant)).Methods("POST")
	// }
	// r.HandleFunc("/admin/applicants/{id}", middleware.AuthMiddleware(handler.GetApplicantByID)).Methods("GET")
	// r.HandleFunc("/admin/applicants/{id}", middleware.AuthMiddleware(handler.UpdateApplicant)).Methods("PUT")
	// r.HandleFunc("/admin/applicants/{id}", middleware.AuthMiddleware(handler.DeleteApplicant)).Methods("DELETE")

	// // Job Management Routes (Optional)
	// r.HandleFunc("/admin/jobs", middleware.AuthMiddleware(handler.GetAllJobs)).Methods("GET")
	// r.HandleFunc("/admin/jobs/{id}", middleware.AuthMiddleware(handler.GetJobByID)).Methods("GET")
	// r.HandleFunc("/admin/jobs/{id}", middleware.AuthMiddleware(handler.DeleteJob)).Methods("DELETE")

	// // Application Management Routes (Optional)
	// r.HandleFunc("/admin/applications", middleware.AuthMiddleware(handler.GetAllApplications)).Methods("GET")
	// r.HandleFunc("/admin/applications/{id}", middleware.AuthMiddleware(handler.GetApplicationByID)).Methods("GET")
	// r.HandleFunc("/admin/applications/{id}", middleware.AuthMiddleware(handler.DeleteApplication)).Methods("DELETE")

}
