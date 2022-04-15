/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package simplexds provides a very simple in-memory XDS server, based on the fake in the grpc package.
// Mainly for testing - the fake is in an internal package.
package uxds

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//adsgrpc "github.com/costinm/grpc-mesh/xds"
	//discoverypb "github.com/costinm/grpc-mesh/xds"
	//lrsgrpc "github.com/costinm/grpc-mesh/xds"
	//lrspb "github.com/costinm/grpc-mesh/xds"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	adsgrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	lrsgrpc "github.com/envoyproxy/go-control-plane/envoy/service/load_stats/v2"
	lrspb "github.com/envoyproxy/go-control-plane/envoy/service/load_stats/v2"
)

const (
	// TODO: Make this a var or a field in the server if there is a need to use a
	// value other than this default.
	defaultChannelBufferSize = 50
	defaultDialTimeout       = 5 * time.Second
)

// Request wraps the request protobuf (xds/LRS) and error received by the
// Server in a call to stream.Recv().
type Request struct {
	Req proto.Message
	Err error
}

// Response wraps the response protobuf (xds/LRS) and error that the Server
// should send out to the client through a call to stream.Send()
type Response struct {
	Resp proto.Message
	Err  error
}

// Server is a micro implementation of xDS and LRS protocols. It listens on the
// same port for both services and exposes a bunch of channels to send/receive
// messages.
// Based on grpc fake XDS server used in tests.
type Server struct {
	// XDSRequestChan is a channel on which received xDS requests are made
	// available to the users of this Server.
	XDSRequestChan *Channel

	// XDSResponseChan is a channel on which the Server accepts xDS responses
	// to be sent to the client.
	XDSResponseChan chan *Response

	// LRSRequestChan is a channel on which received LRS requests are made
	// available to the users of this Server.
	LRSRequestChan *Channel

	// LRSResponseChan is a channel on which the Server accepts the LRS
	// response to be sent to the client.
	LRSResponseChan chan *Response

	// NewConnChan is a channel on which the fake server notifies receipt of new
	// connection attempts. Tests can gate on this event before proceeding to
	// other actions which depend on a connection to the fake server being up.
	NewConnChan *Channel

	// Address is the host:port on which the Server is listening for requests.
	Address string

	// The underlying fake implementation of xDS and LRS.
	xdsS *xdsServer
	lrsS *lrsServer
}

type wrappedListener struct {
	net.Listener
	server *Server
}

func (wl *wrappedListener) Accept() (net.Conn, error) {
	c, err := wl.Listener.Accept()
	if err != nil {
		return nil, err
	}
	wl.server.NewConnChan.Send(struct{}{})
	return c, err
}

// StartServer makes a new Server and gets it to start listening on a local
// port for gRPC requests. The returned cancel function should be invoked by
// the caller upon completion of the test.
func StartServer(addr string, onConnection func(*Connection, *discoverypb.DiscoveryRequest)) (*Server, func(), error) {
	if addr == "" {
		addr = "localhost:0"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, func() {}, fmt.Errorf("net.Listen() failed: %v", err)
	}

	s := &Server{
		XDSRequestChan:  NewChannelWithSize(defaultChannelBufferSize),
		LRSRequestChan:  NewChannelWithSize(defaultChannelBufferSize),
		NewConnChan:     NewChannelWithSize(defaultChannelBufferSize),
		XDSResponseChan: make(chan *Response, defaultChannelBufferSize),
		LRSResponseChan: make(chan *Response, 1), // The server only ever sends one response.
		Address:         lis.Addr().String(),
	}
	s.xdsS = &xdsServer{reqChan: s.XDSRequestChan, respChan: s.XDSResponseChan, clients: map[string]*Connection{}, conCallback: onConnection}
	s.lrsS = &lrsServer{reqChan: s.LRSRequestChan, respChan: s.LRSResponseChan}
	wp := &wrappedListener{
		Listener: lis,
		server:   s,
	}

	server := grpc.NewServer()
	lrsgrpc.RegisterLoadReportingServiceServer(server, s.lrsS)
	adsgrpc.RegisterAggregatedDiscoveryServiceServer(server, s.xdsS)
	go server.Serve(wp)

	return s, func() { server.Stop() }, nil
}

type xdsServer struct {
	reqChan  *Channel
	respChan chan *Response

	// Extra for multi-connection
	// mutex used to modify structs, non-blocking code only.
	mutex sync.RWMutex

	// clients reflect active gRPC channels, for both ADS and MCP.
	// key is Connection.ConID
	clients          map[string]*Connection
	conCallback      func(*Connection, *discoverypb.DiscoveryRequest)
	connectionNumber int
}

func (xdsS *xdsServer) StreamAggregatedResources(s adsgrpc.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	con, f := xdsS.newConnection(s)
	defer f()

	errCh := make(chan error, 2)
	go func() {
		for {
			req, err := s.Recv()
			if xdsS.conCallback != nil {
				if con.firstReq {
					xdsS.conCallback(con, req)
					con.firstReq = false
				}
				con.OnMessage(req, err)
				if err != nil {
					return
				}
			} else {
				if err != nil {
					errCh <- err
					return
				}
				xdsS.reqChan.Send(&Request{req, err})
			}
		}
	}()
	go func() {
		var retErr error
		defer func() {
			errCh <- retErr
		}()

		for {
			select {
			case r := <-con.resChannel:
				if r.Err != nil {
					retErr = r.Err
					return
				}
				if err := s.Send(r.Resp.(*discoverypb.DiscoveryResponse)); err != nil {
					retErr = err
					return
				}
			case r := <-xdsS.respChan:
				if r.Err != nil {
					retErr = r.Err
					return
				}
				if err := s.Send(r.Resp.(*discoverypb.DiscoveryResponse)); err != nil {
					retErr = err
					return
				}
			case <-s.Context().Done():
				retErr = s.Context().Err()
				return
			}
		}
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (xdsS *xdsServer) DeltaAggregatedResources(adsgrpc.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	return status.Error(codes.Unimplemented, "")
}

type lrsServer struct {
	reqChan  *Channel
	respChan chan *Response
}

func (lrsS *lrsServer) StreamLoadStats(s lrsgrpc.LoadReportingService_StreamLoadStatsServer) error {
	req, err := s.Recv()
	lrsS.reqChan.Send(&Request{req, err})
	if err != nil {
		return err
	}

	select {
	case r := <-lrsS.respChan:
		if r.Err != nil {
			return r.Err
		}
		if err := s.Send(r.Resp.(*lrspb.LoadStatsResponse)); err != nil {
			return err
		}
	case <-s.Context().Done():
		return s.Context().Err()
	}

	for {
		req, err := s.Recv()
		lrsS.reqChan.Send(&Request{req, err})
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
