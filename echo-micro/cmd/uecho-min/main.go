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
	"os"
	"os/signal"
	"syscall"

	"github.com/costinm/grpc-mesh/echo-micro/proto"
	"github.com/costinm/grpc-mesh/echo-micro/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(port string) error {
	h := &server.EchoGrpcHandler{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	creds := insecure.NewCredentials()

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	proto.RegisterEchoTestServiceServer(grpcServer, h)

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}

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
//
// 	 ocgrpc adds ~300k
func main() {
	err := Run(":8080")
	if err != nil {
		fmt.Println("Error ", err)
	}
}
