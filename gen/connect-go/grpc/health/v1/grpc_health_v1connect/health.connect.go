// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: grpc/health/v1/health.proto

package grpc_health_v1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// HealthName is the fully-qualified name of the Health service.
	HealthName = "grpc.health.v1.Health"
)

// HealthClient is a client for the grpc.health.v1.Health service.
type HealthClient interface {
	// If the requested service is unknown, the call will fail with status
	// NOT_FOUND.
	Check(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.Response[grpc_health_v1.HealthCheckResponse], error)
	// Performs a watch for the serving status of the requested service.
	// The server will immediately send back a message indicating the current
	// serving status.  It will then subsequently send a new message whenever
	// the service's serving status changes.
	//
	// If the requested service is unknown when the call is received, the
	// server will send a message setting the serving status to
	// SERVICE_UNKNOWN but will *not* terminate the call.  If at some
	// future point, the serving status of the service becomes known, the
	// server will send a new message with the service's serving status.
	//
	// If the call terminates with status UNIMPLEMENTED, then clients
	// should assume this method is not supported and should not retry the
	// call.  If the call terminates with any other status (including OK),
	// clients should retry the call with appropriate exponential backoff.
	Watch(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.ServerStreamForClient[grpc_health_v1.HealthCheckResponse], error)
}

// NewHealthClient constructs a client for the grpc.health.v1.Health service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewHealthClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) HealthClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &healthClient{
		check: connect_go.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
			httpClient,
			baseURL+"/grpc.health.v1.Health/Check",
			opts...,
		),
		watch: connect_go.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
			httpClient,
			baseURL+"/grpc.health.v1.Health/Watch",
			opts...,
		),
	}
}

// healthClient implements HealthClient.
type healthClient struct {
	check *connect_go.Client[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse]
	watch *connect_go.Client[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse]
}

// Check calls grpc.health.v1.Health.Check.
func (c *healthClient) Check(ctx context.Context, req *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.Response[grpc_health_v1.HealthCheckResponse], error) {
	return c.check.CallUnary(ctx, req)
}

// Watch calls grpc.health.v1.Health.Watch.
func (c *healthClient) Watch(ctx context.Context, req *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.ServerStreamForClient[grpc_health_v1.HealthCheckResponse], error) {
	return c.watch.CallServerStream(ctx, req)
}

// HealthHandler is an implementation of the grpc.health.v1.Health service.
type HealthHandler interface {
	// If the requested service is unknown, the call will fail with status
	// NOT_FOUND.
	Check(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.Response[grpc_health_v1.HealthCheckResponse], error)
	// Performs a watch for the serving status of the requested service.
	// The server will immediately send back a message indicating the current
	// serving status.  It will then subsequently send a new message whenever
	// the service's serving status changes.
	//
	// If the requested service is unknown when the call is received, the
	// server will send a message setting the serving status to
	// SERVICE_UNKNOWN but will *not* terminate the call.  If at some
	// future point, the serving status of the service becomes known, the
	// server will send a new message with the service's serving status.
	//
	// If the call terminates with status UNIMPLEMENTED, then clients
	// should assume this method is not supported and should not retry the
	// call.  If the call terminates with any other status (including OK),
	// clients should retry the call with appropriate exponential backoff.
	Watch(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest], *connect_go.ServerStream[grpc_health_v1.HealthCheckResponse]) error
}

// NewHealthHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewHealthHandler(svc HealthHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/grpc.health.v1.Health/Check", connect_go.NewUnaryHandler(
		"/grpc.health.v1.Health/Check",
		svc.Check,
		opts...,
	))
	mux.Handle("/grpc.health.v1.Health/Watch", connect_go.NewServerStreamHandler(
		"/grpc.health.v1.Health/Watch",
		svc.Watch,
		opts...,
	))
	return "/grpc.health.v1.Health/", mux
}

// UnimplementedHealthHandler returns CodeUnimplemented from all methods.
type UnimplementedHealthHandler struct{}

func (UnimplementedHealthHandler) Check(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest]) (*connect_go.Response[grpc_health_v1.HealthCheckResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("grpc.health.v1.Health.Check is not implemented"))
}

func (UnimplementedHealthHandler) Watch(context.Context, *connect_go.Request[grpc_health_v1.HealthCheckRequest], *connect_go.ServerStream[grpc_health_v1.HealthCheckResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("grpc.health.v1.Health.Watch is not implemented"))
}
