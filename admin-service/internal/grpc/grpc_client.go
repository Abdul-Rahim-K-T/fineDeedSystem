package grpc

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	employerpb "fineDeedSystem/proto/employer"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type EmployerClient struct {
// 	client employerpb.EmployerServiceClient
// 	conn   *grpc.ClientConn
// }

// func NewEmployerClient(address string) (*EmployerClient, error) {
// 	conn, err := grpc.Dial(
// 		address,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := employerpb.NewEmployerServiceClient(conn)
// 	return &EmployerClient{client: client, conn: conn}, nil
// }

// func (ec *EmployerClient) GetAllEmployers(ctx context.Context, in *employerpb.Empty) (*employerpb.EmployerList, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	log.Println("Sending request to GetAllEmployers...")
// 	fmt.Println("Sending request to GetAllEmployers...")
// 	log.Printf("Request: %+v\n", in) // Log the request
// 	fmt.Printf("Request: %+v\n", in)

// 	employerList, err := ec.client.GetAllEmployers(ctx, in)
// 	if err != nil {
// 		log.Printf("Error while calling GetAllEmployers: %v\n", err)
// 		return nil, err
// 	}

// 	log.Printf("Response received: %+v\n", employerList) // Log the response
// 	fmt.Printf("Response received: %+v\n", employerList)
// 	return employerList, nil
// }

// func (ec *EmployerClient) Close() error {
// 	return ec.conn.Close()
// }
