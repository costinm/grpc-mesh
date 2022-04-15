package uxds

import (
	"fmt"
	"io"
	"sync"
	"time"

	discoverypb "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	adsgrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/peer"
)

// Connection represents a single endpoint.
// An endpoint typically has 0 or 1 connections - but during restarts and drain it may have >1.
type Connection struct {
	mu sync.RWMutex

	// PeerAddr is the address of the client envoy, from network layer
	PeerAddr string

	NodeID string

	// Time of connection, for debugging
	Connect time.Time

	// ConID is the connection identifier, used as a key in the connection table.
	// Currently based on the node name and a counter.
	ConID string

	// doneChannel will be closed when the client is closed.
	doneChannel chan int

	// Metadata key-value pairs extending the Node identifier
	Metadata map[string]string

	// Watched resources for the connection
	Watched map[string][]string

	NonceSent  map[string]string
	NonceAcked map[string]string

	// Only one can be set.
	SStream adsgrpc.AggregatedDiscoveryService_StreamAggregatedResourcesServer
	CStream adsgrpc.AggregatedDiscoveryService_StreamAggregatedResourcesClient

	active     bool
	resChannel chan *Response
	errChannel chan error

	firstReq bool
}

func (con *Connection) OnMessage(req *discoverypb.DiscoveryRequest, err error) {

	if err != nil {
		if status.Code(err) == codes.Canceled || err == io.EOF {
			con.errChannel <- nil
			return
		}
		con.errChannel <- err
		return
	}
	// TODO

}

func (s *xdsServer) newConnection(stream adsgrpc.AggregatedDiscoveryService_StreamAggregatedResourcesServer) (*Connection, func()) {
	// TODO: add a timer to verify first packet is received
	peerInfo, ok := peer.FromContext(stream.Context())
	peerAddr := "0.0.0.0"
	if ok {
		peerAddr = peerInfo.Addr.String()
	}

	t0 := time.Now()

	con := &Connection{
		Connect:     t0,
		PeerAddr:    peerAddr,
		SStream:     stream,
		NonceSent:   map[string]string{},
		Metadata:    map[string]string{},
		Watched:     map[string][]string{},
		NonceAcked:  map[string]string{},
		doneChannel: make(chan int, 2),
		resChannel:  make(chan *Response, 2),
		errChannel:  make(chan error, 2),
		firstReq:    true,
	}
	return con, func() {
		if con.firstReq {
			return // didn't get first req, not added
		}
		close(con.resChannel)
		close(con.doneChannel)
		s.mutex.Lock()
		delete(s.clients, con.ConID)
		s.mutex.Unlock()

	}
}

func (fx *xdsServer) SendAll(r *discoverypb.DiscoveryResponse) {
	for _, con := range fx.clients {
		// TODO: only if watching our resource type

		r.Nonce = fmt.Sprintf("%v", time.Now())
		con.NonceSent[r.TypeUrl] = r.Nonce
		con.resChannel <- &Response{Resp: r}
		// Not safe to call from 2 threads: con.SStream.Send(r)
	}
}

func (fx *xdsServer) Send(con *Connection, r *discoverypb.DiscoveryResponse) {
	r.Nonce = fmt.Sprintf("%v", time.Now())
	con.NonceSent[r.TypeUrl] = r.Nonce
	con.resChannel <- &Response{Resp: r}
}
