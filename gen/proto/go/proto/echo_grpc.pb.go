// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/echo.proto

package proto

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

// EchoTestServiceClient is the client API for EchoTestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EchoTestServiceClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	ForwardEcho(ctx context.Context, in *ForwardEchoRequest, opts ...grpc.CallOption) (*ForwardEchoResponse, error)
}

type echoTestServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEchoTestServiceClient(cc grpc.ClientConnInterface) EchoTestServiceClient {
	return &echoTestServiceClient{cc}
}

func (c *echoTestServiceClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/proto.EchoTestService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoTestServiceClient) ForwardEcho(ctx context.Context, in *ForwardEchoRequest, opts ...grpc.CallOption) (*ForwardEchoResponse, error) {
	out := new(ForwardEchoResponse)
	err := c.cc.Invoke(ctx, "/proto.EchoTestService/ForwardEcho", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoTestServiceServer is the server API for EchoTestService service.
// All implementations should embed UnimplementedEchoTestServiceServer
// for forward compatibility
type EchoTestServiceServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	ForwardEcho(context.Context, *ForwardEchoRequest) (*ForwardEchoResponse, error)
}

// UnimplementedEchoTestServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEchoTestServiceServer struct {
}

func (UnimplementedEchoTestServiceServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedEchoTestServiceServer) ForwardEcho(context.Context, *ForwardEchoRequest) (*ForwardEchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardEcho not implemented")
}

// UnsafeEchoTestServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EchoTestServiceServer will
// result in compilation errors.
type UnsafeEchoTestServiceServer interface {
	mustEmbedUnimplementedEchoTestServiceServer()
}

func RegisterEchoTestServiceServer(s grpc.ServiceRegistrar, srv EchoTestServiceServer) {
	s.RegisterService(&EchoTestService_ServiceDesc, srv)
}

func _EchoTestService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoTestServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EchoTestService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoTestServiceServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoTestService_ForwardEcho_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardEchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoTestServiceServer).ForwardEcho(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EchoTestService/ForwardEcho",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoTestServiceServer).ForwardEcho(ctx, req.(*ForwardEchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EchoTestService_ServiceDesc is the grpc.ServiceDesc for EchoTestService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EchoTestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EchoTestService",
	HandlerType: (*EchoTestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoTestService_Echo_Handler,
		},
		{
			MethodName: "ForwardEcho",
			Handler:    _EchoTestService_ForwardEcho_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/echo.proto",
}
