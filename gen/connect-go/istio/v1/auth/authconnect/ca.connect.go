// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: istio/v1/auth/ca.proto

package authconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	auth "github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
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
	// IstioCertificateServiceName is the fully-qualified name of the IstioCertificateService service.
	IstioCertificateServiceName = "istio.v1.auth.IstioCertificateService"
)

// IstioCertificateServiceClient is a client for the istio.v1.auth.IstioCertificateService service.
type IstioCertificateServiceClient interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(context.Context, *connect_go.Request[auth.IstioCertificateRequest]) (*connect_go.Response[auth.IstioCertificateResponse], error)
}

// NewIstioCertificateServiceClient constructs a client for the
// istio.v1.auth.IstioCertificateService service. By default, it uses the Connect protocol with the
// binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use the
// gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewIstioCertificateServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) IstioCertificateServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &istioCertificateServiceClient{
		createCertificate: connect_go.NewClient[auth.IstioCertificateRequest, auth.IstioCertificateResponse](
			httpClient,
			baseURL+"/istio.v1.auth.IstioCertificateService/CreateCertificate",
			opts...,
		),
	}
}

// istioCertificateServiceClient implements IstioCertificateServiceClient.
type istioCertificateServiceClient struct {
	createCertificate *connect_go.Client[auth.IstioCertificateRequest, auth.IstioCertificateResponse]
}

// CreateCertificate calls istio.v1.auth.IstioCertificateService.CreateCertificate.
func (c *istioCertificateServiceClient) CreateCertificate(ctx context.Context, req *connect_go.Request[auth.IstioCertificateRequest]) (*connect_go.Response[auth.IstioCertificateResponse], error) {
	return c.createCertificate.CallUnary(ctx, req)
}

// IstioCertificateServiceHandler is an implementation of the istio.v1.auth.IstioCertificateService
// service.
type IstioCertificateServiceHandler interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(context.Context, *connect_go.Request[auth.IstioCertificateRequest]) (*connect_go.Response[auth.IstioCertificateResponse], error)
}

// NewIstioCertificateServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewIstioCertificateServiceHandler(svc IstioCertificateServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/istio.v1.auth.IstioCertificateService/CreateCertificate", connect_go.NewUnaryHandler(
		"/istio.v1.auth.IstioCertificateService/CreateCertificate",
		svc.CreateCertificate,
		opts...,
	))
	return "/istio.v1.auth.IstioCertificateService/", mux
}

// UnimplementedIstioCertificateServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedIstioCertificateServiceHandler struct{}

func (UnimplementedIstioCertificateServiceHandler) CreateCertificate(context.Context, *connect_go.Request[auth.IstioCertificateRequest]) (*connect_go.Response[auth.IstioCertificateResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("istio.v1.auth.IstioCertificateService.CreateCertificate is not implemented"))
}
