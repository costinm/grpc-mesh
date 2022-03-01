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
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/costinm/grpc-mesh/echo-micro/proto"
	"github.com/costinm/grpc-mesh/echo-micro/server"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/health"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

var log = grpclog.Component("echo")

func Run(port string) {
	h := &server.EchoGrpcHandler{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("net.Listen(tcp)", port, err)
	}

	creds := insecure.NewCredentials()

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
	}

	// Special handling for startup without env variable set

	// Generate the bootstrap if the file is missing ( injection-less )
	// using cloudrun-mesh auto-detection code

	// Generate certs if missing

	grpcServer := grpc.NewServer(grpcOptions...)
	proto.RegisterEchoTestServiceServer(grpcServer, h)

	// add the standard grpc health check
	healthServer := health.NewServer()
	grpcHealth.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)
	cleanup, err := admin.Register(grpcServer)
	if err != nil {
		log.Info("Failed to register admin", "error", err)
	}

	go func() {
		err = grpcServer.Serve(lis)
		log.Info("grpcServer done")
	}()

	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	cleanup()
	// TODO: lame duck, etc
}

func main() {
	Run(":8080")
}
