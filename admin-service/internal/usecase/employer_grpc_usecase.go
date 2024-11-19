package usecase

import (
	"context"
	"fineDeedSystem/admin-service/proto/fineDeedSystem/proto/shared"
	"fmt"
	"log"
)

type EmployerGrpcUsecase struct {
	grpcClient shared.EmployerServiceClient
}

func NewEmployerGrpcUsecase(grpcClient shared.EmployerServiceClient) *EmployerGrpcUsecase {
	return &EmployerGrpcUsecase{
		grpcClient: grpcClient,
	}
}

func (u *EmployerGrpcUsecase) GrpcCreateEmployer(ctx context.Context, employer *shared.Employer) (*shared.Employer, error) {
	fmt.Println("GrpcCreateEmployer at usecase/employer_grpc_usecase.go file")
	createRequest := &shared.CreateEmployerRequest{
		Employer: employer,
	}

	// Assuming you have a GrpcCreateEmployer RPC in your employer-service
	resp, err := u.grpcClient.CreateEmployer(ctx, createRequest)
	if err != nil {
		log.Println("Error calling GrpcCreateEmployer:", err)
		return nil, err
	}
	return resp.Employer, nil
}

// func (u *EmployerGrpcUsecase) GetAllEmployers(ctx context.Context) ([]*shared.Employer, error) {
// 	fmt.Println("GetAllEmployers at usecase/employer_grpc_usecase.go file")
// 	// Assuming you have a GetAllEmployers RPC in your employer-service
// 	employers, err := u.grpcClient.ListEmployers(ctx, &shared.ListEmployersRequest{})
// 	if err != nil {
// 		log.Printf("Error from gRPC client: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("Received from gRPC client: %+v", employers)
// 	return employers.Employers, nil
// }
