package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/costinm/grpc-mesh/gen/connect/go/envoy/service/discovery/v3/v3connect"
	"github.com/costinm/grpc-mesh/gen/connect/go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/pkg/echo"
	"github.com/costinm/grpc-mesh/pkg/goh2"
	"github.com/costinm/grpc-mesh/pkg/istioca"
	"github.com/costinm/grpc-mesh/pkg/xds"
	"github.com/costinm/meshauth"
	meshca "github.com/costinm/meshauth/pkg/ca"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	ctx := context.Background()
	ug, _ := meshauth.FromEnv(ctx, nil, "umeshd")

	// Attempt to load roots from filesystem
	ca := meshca.NewCA()

	mux := http.NewServeMux()

	mux.Handle(authconnect.NewIstioCertificateServiceHandler(&istioca.IstioCA{CA: ca, Auth: ug}))

	xdss := &xds.XDSServer{}
	echos := &echo.Echo{}

	mux.Handle(v3connect.NewAggregatedDiscoveryServiceHandler(xdss))

	echos.RegisterMux(mux, "")

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	// Init h2 server module.
	h2 := &goh2.H2{}
	h2.Init(ctx)

	h2.Start(ctx)
	// Init h2 server
	h2s := &http2.Server{}
	h2ch := h2c.NewHandler(mux, h2s)
	http.Serve(l, h2ch)
}
