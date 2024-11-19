package grpc

import (
	"context"
	"fineDeedSystem/employer-service/internal/models"
	"fineDeedSystem/employer-service/internal/usecase"
	"fineDeedSystem/employer-service/proto/fineDeedSystem/proto/shared"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EmployerGrpcHandler struct {
	usecase *usecase.EmployerUsecase
	shared.UnimplementedEmployerServiceServer
}

func NewEmployerGrpcHandler(usecase *usecase.EmployerUsecase) *EmployerGrpcHandler {
	return &EmployerGrpcHandler{usecase: usecase}
}

func (e *EmployerGrpcHandler) DeleteEmployer(ctx context.Context, req *shared.DeleteEmployerRequest) (*shared.DeleteEmployerResponse, error) {
	// Call the usecase method to delete the employer
	err := e.usecase.DeleteEmployer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete employer: %v", err)
	}

	// Return a success message
	return &shared.DeleteEmployerResponse{
		Message: "Employer deleted successfully",
	}, nil
}

func (e *EmployerGrpcHandler) UpdateEmployer(ctx context.Context, req *shared.UpdateEmployerRequest) (*shared.UpdateEmployerResponse, error) {
	// // Fetch the existing employer details
	// existingEmployer, err := e.usecase.GetEmployerByID(ctx, uint32(req.Employer.Id))
	// if err != nil {
	// 	return nil, err
	// }
	// if existingEmployer == nil {
	// 	return nil, errors.New("employer not found")
	// }

	// // Compare and update fields
	// updated := false
	// if req.Employer.Name != "" && req.Employer.Name != existingEmployer.Name {
	// 	existingEmployer.Name = req.Employer.Name
	// 	updated = true
	// }
	// if req.Employer.Email != "" && req.Employer.Email != existingEmployer.Email {
	// 	existingEmployer.Email = req.Employer.Email
	// 	updated = true
	// }
	// if req.Employer.Phone != "" && req.Employer.Phone != existingEmployer.Phone {
	// 	existingEmployer.Phone = req.Employer.Phone
	// 	updated = true
	// }
	// if req.Employer.CompanyName != "" && req.Employer.CompanyName != existingEmployer.CompanyName {
	// 	existingEmployer.CompanyName = req.Employer.CompanyName
	// 	updated = true
	// }
	// if req.Employer.Password != "" && req.Employer.Password != existingEmployer.Password {
	// 	existingEmployer.Password = req.Employer.Password
	// 	updated = true
	// }

	// if !updated {
	// 	return nil, errors.New("no fields were changed")
	// }

	// Extract employer details from the request
	updatedEmployer := &models.Employer{
		ID:          uint(req.Employer.Id),
		Name:        req.Employer.Name,
		Email:       req.Employer.Email,
		Phone:       req.Employer.Phone,
		CompanyName: req.Employer.CompanyName,
		Password:    req.Employer.Password,
	}

	// call the usecase method to update the employer
	employer, err := e.usecase.UpdateEmployer(ctx, updatedEmployer)
	if err != nil {
		// Check if the errror is for no fields changed and return a specific response
		if err.Error() == "no fields are changed" {
			return nil, status.Errorf(codes.InvalidArgument, "no fields are changed")
		}
		return nil, err
	}

	// Return the updated employer in the response
	return &shared.UpdateEmployerResponse{
		Employer: &shared.Employer{
			Id:          uint32(employer.ID),
			Name:        employer.Name,
			Email:       employer.Email,
			Phone:       employer.Phone,
			CompanyName: employer.CompanyName,
		},
	}, nil
}

func (e *EmployerGrpcHandler) GetEmployerByID(ctx context.Context, req *shared.GetEmployerByIdRequest) (*shared.GetEmployerByIdResponse, error) {
	// Retrieve employer ID from the request
	employerID := req.Id

	// Call the usecase method to get the employer details
	employer, err := e.usecase.GetEmployerByID(ctx, employerID)
	if err != nil {
		return nil, err
	}

	// Debugging log
	log.Printf("Fetched employer details: ID=%d, Name=%s, Email=%s, Phone=%s, CompanyName=%s",
		employer.ID, employer.Name, employer.Email, employer.Phone, employer.CompanyName)

	// Return the response with employer details
	return &shared.GetEmployerByIdResponse{
		Employer: &shared.Employer{
			Id:          uint32(employer.ID),
			Name:        employer.Name,
			Email:       employer.Email,
			Phone:       employer.Phone,
			CompanyName: employer.CompanyName,
		},
	}, nil
}

func (e *EmployerGrpcHandler) CreateEmployer(ctx context.Context, req *shared.CreateEmployerRequest) (*shared.CreateEmployerResponse, error) {
	employer := &models.Employer{
		Name:        req.Employer.Name,
		Email:       req.Employer.Email,
		Phone:       req.Employer.Phone,
		CompanyName: req.Employer.CompanyName,
		Password:    req.Employer.Password,
	}

	// Debugging log
	log.Printf("Received employer details: Name=%s, Email=%s, Phone=%s,CompanyName=%s, Password=%s",
		employer.Name, employer.Email, employer.Phone, employer.CompanyName, employer.Password)

	newEmployer, err := e.usecase.CreateEmployerLogic(employer)
	if err != nil {
		return nil, err
	}

	return &shared.CreateEmployerResponse{
		Employer: &shared.Employer{
			Id:          uint32(newEmployer.ID),
			Name:        newEmployer.Name,
			Email:       newEmployer.Email,
			Phone:       newEmployer.Phone,
			CompanyName: newEmployer.CompanyName,
			Password:    newEmployer.Password,
		},
	}, nil
}

func (e *EmployerGrpcHandler) ListEmployers(ctx context.Context, req *shared.ListEmployersRequest) (*shared.ListEmployersResponse, error) {
	employers, err := e.usecase.GetAllEmployers()
	if err != nil {
		return nil, err
	}

	var employerList []*shared.Employer
	for _, emp := range employers {
		employerList = append(employerList, &shared.Employer{
			Id:          uint32(emp.ID),
			Name:        emp.Name,
			Email:       emp.Email,
			Phone:       emp.Phone,
			CompanyName: emp.CompanyName,
		})
	}
	return &shared.ListEmployersResponse{Employers: employerList}, nil
}
