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
	"os"
	"os/signal"
	"syscall"

	"github.com/costinm/grpc-mesh/echo-min/common"
	// To install the xds resolvers and balancers.
	_ "google.golang.org/grpc/xds"

	"github.com/costinm/grpc-mesh/echo-min/server"
	"istio.io/pkg/log"
)

var (
	httpPorts        []int
	grpcPorts        []int
	tcpPorts         []int
	tlsPorts         []int
	instanceIPPorts  []int
	localhostIPPorts []int
	serverFirstPorts []int
	xdsGRPCServers   []int
	metricsPort      int
	uds              string
	version          string
	cluster          string
	crt              string
	key              string
	istioVersion     string
	disableALPN      bool

	loggingOptions = log.DefaultOptions()
)

func Run(port string) {
	instanceIPByPort := map[int]struct{}{}
	for _, p := range instanceIPPorts {
		instanceIPByPort[p] = struct{}{}
	}
	localhostIPByPort := map[int]struct{}{}
	for _, p := range localhostIPPorts {
		localhostIPByPort[p] = struct{}{}
	}

	s := server.New(server.Config{
		Ports:                 []*common.Port{},
		Metrics:               metricsPort,
		BindIPPortsMap:        instanceIPByPort,
		BindLocalhostPortsMap: localhostIPByPort,
		TLSCert:               crt,
		TLSKey:                key,
		Version:               version,
		Cluster:               cluster,
		IstioVersion:          istioVersion,
		UDSServer:             uds,
		DisableALPN:           disableALPN,
	})

	if err := s.Start(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	defer func() {
		_ = s.Close()
	}()

	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

}

func main() {
	Run(":8080")
}
