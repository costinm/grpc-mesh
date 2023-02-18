package main

import (
	"context"
	"flag"
	"github.com/costinm/grpc-mesh/gen/connect-go/envoy/service/load_stats/v2/v2connect"

	"github.com/costinm/grpc-mesh/ptrds/uxds"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect-go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	IstiodAddr = flag.String("xds", "http://localhost:15010", "Istio address")
	ns         = flag.String("n", "istio-system", "")
	ip         = flag.String("ip", "127.1.2.3", "")
	pod        = flag.String("pod", "test-abc-def", "")
)

func main() {
	mux := http.NewServeMux()
	// Mount some handlers here.
	server := &http.Server{
		Addr:    ":http",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		// Don't forget timeouts!
	}

	mux.Handle(v2connect.NewLoadReportingServiceHandler(&uxds.LRS{}))
	//mux.Handle(v3connect.NewAggregatedDiscoveryServiceHandler(&uxds.XDS{})//uxds.NewXDS(&uxds.XDSConfig{})))

	server.ListenAndServe()
}

func GetCert() {
	caclient := authconnect.NewIstioCertificateServiceClient(http.DefaultClient, "")
	req := connect.NewRequest(&auth.IstioCertificateRequest{})

	certificate, err := caclient.CreateCertificate(context.Background(), req)
	if err != nil {
		return
	}
	log.Println(certificate.Msg)
}
