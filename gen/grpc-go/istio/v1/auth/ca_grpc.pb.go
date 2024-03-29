// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: istio/v1/auth/ca.proto

package auth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IstioCertificateServiceClient is the client API for IstioCertificateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IstioCertificateServiceClient interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(ctx context.Context, in *IstioCertificateRequest, opts ...grpc.CallOption) (*IstioCertificateResponse, error)
}

type istioCertificateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIstioCertificateServiceClient(cc grpc.ClientConnInterface) IstioCertificateServiceClient {
	return &istioCertificateServiceClient{cc}
}

func (c *istioCertificateServiceClient) CreateCertificate(ctx context.Context, in *IstioCertificateRequest, opts ...grpc.CallOption) (*IstioCertificateResponse, error) {
	out := new(IstioCertificateResponse)
	err := c.cc.Invoke(ctx, "/istio.v1.auth.IstioCertificateService/CreateCertificate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IstioCertificateServiceServer is the server API for IstioCertificateService service.
// All implementations should embed UnimplementedIstioCertificateServiceServer
// for forward compatibility
type IstioCertificateServiceServer interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(context.Context, *IstioCertificateRequest) (*IstioCertificateResponse, error)
}

// UnimplementedIstioCertificateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedIstioCertificateServiceServer struct {
}

func (UnimplementedIstioCertificateServiceServer) CreateCertificate(context.Context, *IstioCertificateRequest) (*IstioCertificateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCertificate not implemented")
}

// UnsafeIstioCertificateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IstioCertificateServiceServer will
// result in compilation errors.
type UnsafeIstioCertificateServiceServer interface {
	mustEmbedUnimplementedIstioCertificateServiceServer()
}

func RegisterIstioCertificateServiceServer(s grpc.ServiceRegistrar, srv IstioCertificateServiceServer) {
	s.RegisterService(&IstioCertificateService_ServiceDesc, srv)
}

func _IstioCertificateService_CreateCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IstioCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IstioCertificateServiceServer).CreateCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/istio.v1.auth.IstioCertificateService/CreateCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IstioCertificateServiceServer).CreateCertificate(ctx, req.(*IstioCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IstioCertificateService_ServiceDesc is the grpc.ServiceDesc for IstioCertificateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IstioCertificateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "istio.v1.auth.IstioCertificateService",
	HandlerType: (*IstioCertificateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCertificate",
			Handler:    _IstioCertificateService_CreateCertificate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "istio/v1/auth/ca.proto",
}
