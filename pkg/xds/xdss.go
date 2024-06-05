package xds

import (
	"context"
	"log"
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/proto/go/xds"
)

type Config struct {
	// Will maintain connections to each XDS server in the config
	Servers map[string]*XDSConfig
}

// XDSServer is a uXDS server.
// Will hold a number of client connections ( config or dynamic ), plus data loaded
// programmatically.
type XDSServer struct {
	m sync.Mutex

	// Active client connections for full - used to send pushes.
	XDSc map[string]*XDS
}

func NewXDSServer(cfg *Config) *XDSServer {

	return &XDSServer{}
}

func (X *XDSServer) StreamAggregatedResources(ctx context.Context, c *connect.BidiStream[xds.DiscoveryRequest, xds.DiscoveryResponse]) error {
	log.Println(c.Peer(), c.RequestHeader())
	req, err := c.Receive()
	if err != nil {
		return err
	}
	log.Println(req)
	return nil
}

func (X *XDSServer) DeltaAggregatedResources(ctx context.Context, c *connect.BidiStream[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse]) error {
	log.Println(c.Peer(), c.RequestHeader())
	req, err := c.Receive()
	if err != nil {
		return err
	}
	log.Println(req)
	return nil
}

