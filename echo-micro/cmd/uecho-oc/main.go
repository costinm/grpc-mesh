// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/costinm/grpc-mesh/echo-micro/proto"
	"github.com/costinm/grpc-mesh/echo-micro/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	"log"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/runmetrics"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"

	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/reflection"
)

func init() {
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Println("Failed to register ocgrpc server views: %v", err)
	}
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Println("Failed to register ocgrpc server views: %v", err)
	}

	// Similar with pilot-agent
	registry := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWithPrefix("sshca_",
		prometheus.Registerer(registry))

	wrapped.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	wrapped.MustRegister(collectors.NewGoCollector())

	//promRegistry = registry
	// go collector metrics collide with other metrics.
	exporter, err := ocprom.NewExporter(ocprom.Options{Registry: registry,
		Registerer: wrapped})
	if err != nil {
		log.Fatalf("could not setup exporter: %v", err)
	}
	view.RegisterExporter(exporter)

}

func Run(port string) error {
	h := &server.EchoGrpcHandler{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	creds := insecure.NewCredentials()

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}

	err = runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
		Prefix:       "echo/",
	})
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	proto.RegisterEchoTestServiceServer(grpcServer, h)
	admin.Register(grpcServer)
	reflection.Register(grpcServer)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	// Status
	mux := &http.ServeMux{}
	zpages.Handle(mux, "/debug")

	http.ListenAndServe("127.0.0.1:8081", mux)
	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	return nil
}

// Most minimal gRPC based server, for estimating binary size overhead.
//
// - 0.8M for a min go program
// - 4.7M for an echo using HTTP.
// - 9M - this server, only plain gRPC
// - 20M - same app, but proxyless gRPC
// - 22M - plus opencensus, prom, zpages, reflection
func main() {
	err := Run(":8080")
	if err != nil {
		fmt.Println("Error ", err)
	}
}
