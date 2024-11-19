// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: proto/shared.proto

package shared

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EmployerService_ListEmployers_FullMethodName  = "/shared.EmployerService/ListEmployers"
	EmployerService_CreateEmployer_FullMethodName = "/shared.EmployerService/CreateEmployer"
)

// EmployerServiceClient is the client API for EmployerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The employer service definition.
type EmployerServiceClient interface {
	ListEmployers(ctx context.Context, in *ListEmployersRequest, opts ...grpc.CallOption) (*ListEmployersResponse, error)
	CreateEmployer(ctx context.Context, in *CreateEmployerRequest, opts ...grpc.CallOption) (*CreateEmployerResponse, error)
}

type employerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmployerServiceClient(cc grpc.ClientConnInterface) EmployerServiceClient {
	return &employerServiceClient{cc}
}

func (c *employerServiceClient) ListEmployers(ctx context.Context, in *ListEmployersRequest, opts ...grpc.CallOption) (*ListEmployersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEmployersResponse)
	err := c.cc.Invoke(ctx, EmployerService_ListEmployers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *employerServiceClient) CreateEmployer(ctx context.Context, in *CreateEmployerRequest, opts ...grpc.CallOption) (*CreateEmployerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateEmployerResponse)
	err := c.cc.Invoke(ctx, EmployerService_CreateEmployer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmployerServiceServer is the server API for EmployerService service.
// All implementations must embed UnimplementedEmployerServiceServer
// for forward compatibility.
//
// The employer service definition.
type EmployerServiceServer interface {
	ListEmployers(context.Context, *ListEmployersRequest) (*ListEmployersResponse, error)
	CreateEmployer(context.Context, *CreateEmployerRequest) (*CreateEmployerResponse, error)
	mustEmbedUnimplementedEmployerServiceServer()
}

// UnimplementedEmployerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEmployerServiceServer struct{}

func (UnimplementedEmployerServiceServer) ListEmployers(context.Context, *ListEmployersRequest) (*ListEmployersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEmployers not implemented")
}
func (UnimplementedEmployerServiceServer) CreateEmployer(context.Context, *CreateEmployerRequest) (*CreateEmployerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEmployer not implemented")
}
func (UnimplementedEmployerServiceServer) mustEmbedUnimplementedEmployerServiceServer() {}
func (UnimplementedEmployerServiceServer) testEmbeddedByValue()                         {}

// UnsafeEmployerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmployerServiceServer will
// result in compilation errors.
type UnsafeEmployerServiceServer interface {
	mustEmbedUnimplementedEmployerServiceServer()
}

func RegisterEmployerServiceServer(s grpc.ServiceRegistrar, srv EmployerServiceServer) {
	// If the following call pancis, it indicates UnimplementedEmployerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EmployerService_ServiceDesc, srv)
}

func _EmployerService_ListEmployers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEmployersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployerServiceServer).ListEmployers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmployerService_ListEmployers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployerServiceServer).ListEmployers(ctx, req.(*ListEmployersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmployerService_CreateEmployer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEmployerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployerServiceServer).CreateEmployer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmployerService_CreateEmployer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployerServiceServer).CreateEmployer(ctx, req.(*CreateEmployerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmployerService_ServiceDesc is the grpc.ServiceDesc for EmployerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmployerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shared.EmployerService",
	HandlerType: (*EmployerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListEmployers",
			Handler:    _EmployerService_ListEmployers_Handler,
		},
		{
			MethodName: "CreateEmployer",
			Handler:    _EmployerService_CreateEmployer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shared.proto",
}