package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect-go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
	"github.com/costinm/grpc-mesh/gen/proto/go/xds"

	"github.com/costinm/grpc-mesh/gen/connect-go/envoy/service/discovery/v3/v3connect"
)

var (
	IstidAddr = flag.String("xds", "http://localhost:15010", "Istio address")
)

// Client for ZTunnel PTR-DS, using connect-go implementation for evaluation.
func main() {
	client := v3connect.NewAggregatedDiscoveryServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := &xds.DeltaDiscoveryRequest{}
	//req.Header().Set("Some-Header", "hello from connect")
	ctx := context.Background()
	ds := client.DeltaAggregatedResources(ctx)
	err := ds.Send(req)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := ds.Receive()
	log.Println(res)
	//	log.Println(res.Header().Get("Some-Other-Header"))
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
