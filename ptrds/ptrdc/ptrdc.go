package main

import (
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"

	"github.com/costinm/grpc-mesh/gen/connect-go/envoy/service/discovery/v3/v3connect"
)

// Client for ZTunnel PTR-DS, using connect-go implementation for evaluation.
func main() {
	client := v3connect.NewAggregatedDiscoveryServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&xds.Delta{
		Number: 42,
	})
	//req.Header().Set("Some-Header", "hello from connect")
	res, err := client.Ping(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	log.Println(res.Header().Get("Some-Other-Header"))
}
