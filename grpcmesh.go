package grpc_mesh

import (
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/binarylog"
)

// Helper for setting up proxyless GRPC without an agent or envoy.
// Based on cloud-run-mesh, which is used to auto-configure. The defaults
// are based on environment and detected settings.

type MeshConfig struct {
	// BasePort is used to configure the ports. Should be set to 15000 when running in Istio environment.
	// If set to 0, ports will be allocated dynamically.
	// GRPC will use:
	BasePort int

	// If EnvoySidecar is not present, 15020 is the expected port for Istio.
	TelemetryPort int

	// Should be 15021 for Istio.
	HealthPort int

	// 15000 port for Istio without envoy.
	AdminPort int

	// 15011 (old Istio port using mesh identity)
	GRPCMTLSPort int

	// GRPC port to use in 'secure network' mode, when the low-level network provides adequate security.
	//
	Port int

	BinaryLog binarylog.Sink

	Verbose int
}

type GRPCMesh struct {
}

// GRPCServer is the interface implemented by both grpc
type GRPCServer interface {
	RegisterService(*grpc.ServiceDesc, interface{})
	Serve(net.Listener) error
	Stop()
	GracefulStop()
	GetServiceInfo() map[string]grpc.ServiceInfo
}

func MeshSetup(cfg *MeshConfig) (srv *GRPCServer, cleanup func(), err error) {
	port := os.Getenv("PORT")
	var portInt int
	if port == "15009" {
		// Operating in CloudRun or similar environment, in tunnel mode - '--port 15009'.
		// Setup HBone adapter, which terminates MTLS and forwards to non-gRPC ports.
		//
	} else if port == "15008" {
		// Dedicated HBONE mTLS port. Setup the tunnel, GRPC is internally forwarded.
	} else if port != "" {
		portInt, err = strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return
	}

	return
}
