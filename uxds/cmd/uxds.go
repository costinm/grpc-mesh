package main

import (
	"log"
	uxds "simplexds"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
)

// Use the uXDS server with files.
func main() {
	srv, stop, err := uxds.StartServer(":15010",
		func(connection *uxds.Connection, request *envoy_api_v2.DiscoveryRequest) {
			log.Println("Received connection ", request)
		})
	if err != nil {
		panic(err)
	}
	defer stop()
	log.Println("Listen on ", srv.Address)
	select {}
}
