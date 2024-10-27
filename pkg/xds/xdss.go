package xds

import (
	"context"
	"log"
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect/go/envoy/service/discovery/v3/v3connect"
	"github.com/costinm/grpc-mesh/gen/proto/go/xds"
)

type Config struct {
	// Will maintain connections to each XDS server in the config
}

type XDSSC struct {
	c *connect.BidiStream[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse]
}

// XDSServer is a uXDS server.
// Will hold a number of client connections ( config or dynamic ), plus data loaded
// programmatically.
type XDSServer struct {
	m sync.Mutex

	// Active client connections for full - used to send pushes.
	XDSc map[string]*XDSSC

	v3connect.UnimplementedAggregatedDiscoveryServiceHandler
}

var _ v3connect.AggregatedDiscoveryServiceHandler = &XDSServer{}

func NewXDSServer(cfg *Config) *XDSServer {

	return &XDSServer{}
}

func (X *XDSServer) StreamAggregatedResources(ctx context.Context, c *connect.BidiStream[xds.DiscoveryRequest, xds.DiscoveryResponse]) error {
	log.Println(c.Peer(), c.RequestHeader())

	for {
		req, err := c.Receive()
		if err != nil {
			return err
		}
		log.Println(req)
	}
	return nil
}

func (X *XDSServer) DeltaAggregatedResources(ctx context.Context, c *connect.BidiStream[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse]) error {

	var n *xds.Node

	for {
		req, err := c.Receive()
		if err != nil {
			if n != nil {
				delete(X.XDSc, n.Id)
			}
			return err
		}
		if n == nil {
			n = req.Node
			log.Println(c.Peer(), c.RequestHeader(), n)
			X.XDSc[n.Id] = &XDSSC{c: c}
		}

		log.Println(n.Id, req.TypeUrl, req.InitialResourceVersions, req.ResourceNamesSubscribe,
			req.ResourceNamesUnsubscribe, req.ResponseNonce, req.ErrorDetail)

	}
	return nil
}

