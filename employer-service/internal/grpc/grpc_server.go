package grpc

// import (
// 	"context"
// 	"log"
// 	"net"

// 	"fineDeedSystem/employer-service/internal/usecase"
// 	pb "fineDeedSystem/proto/employer"

// 	"google.golang.org/grpc"
// )

// type EmployerServer struct {
// 	pb.UnimplementedEmployerServiceServer
// 	usecase *usecase.EmployerUsecase
// }

// func NewEmployerServer(usecase *usecase.EmployerUsecase) *EmployerServer {
// 	return &EmployerServer{usecase: usecase}
// }

// func (s *EmployerServer) ListEmployers(ctx context.Context, req *pb.ListEmployersRequest) (*pb.ListEmployersResponse, error) {
// 	log.Println("Received request to ListEmployers")

// 	employers, err := s.usecase.GetAllEmployers()
// 	if err != nil {
// 		log.Printf("Error while fetching employers: %v\n", err)
// 		return nil, err
// 	}

// 	employerList := &pb.ListEmployersResponse{}
// 	for _, employer := range employers {
// 		employerList.Employers = append(employerList.Employers, &pb.Employer{
// 			Id:          uint32(employer.ID),
// 			Name:        employer.Name,
// 			Email:       employer.Email,
// 			Phone:       employer.Phone,
// 			CompanyName: employer.CompanyName,
// 		})
// 	}

// 	log.Printf("Responding with employers: %+v\n", employerList)
// 	return employerList, nil
// }

// func StartGRPCServer(port string, usecase *usecase.EmployerUsecase) {
// 	lis, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterEmployerServiceServer(grpcServer, NewEmployerServer(usecase))
// 	log.Printf("gRPC server listening on %s", port)
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
