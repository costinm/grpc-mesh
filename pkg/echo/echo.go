package echo

import (
	"context"
	"log"
	"net/http"
	"time"

	connect_go "github.com/bufbuild/connect-go"
	istioechoconnect "github.com/costinm/grpc-mesh/gen/connect/go/proto/protoconnect"

	"github.com/costinm/grpc-mesh/gen/connect-go/fgrpc/fgrpcconnect"
	"github.com/costinm/grpc-mesh/gen/connect-go/proto/protoconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/fgrpc"
	istioecho "github.com/costinm/grpc-mesh/gen/proto/go/proto"
)

type Echo struct {
	protoconnect.UnimplementedEchoTestServiceHandler
	fgrpcconnect.UnimplementedPingServerHandler
}

func (echos *Echo) RegisterMux(mux *http.ServeMux, prefix string) {
	mux.Handle(istioechoconnect.NewEchoTestServiceHandler(echos))
	// Fortio ping interface
	mux.Handle(fgrpcconnect.NewPingServerHandler(echos))
}

func (*Echo) Echo(ctx context.Context, req *connect_go.Request[istioecho.EchoRequest]) (*connect_go.Response[istioecho.EchoResponse], error) {
	log.Println(req.Header(), req.Peer(), req.Spec(), req.Msg)

	res := connect_go.NewResponse(&istioecho.EchoResponse{
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func (*Echo) Ping(ctx context.Context, req *connect_go.Request[fgrpc.PingMessage]) (*connect_go.Response[fgrpc.PingMessage], error) {

	if req.Msg.DelayNanos > 0 {
		time.Sleep(time.Duration(req.Msg.DelayNanos))
	}

	return connect_go.NewResponse(&fgrpc.PingMessage{
		Seq:        req.Msg.Seq,
		Ts:         req.Msg.Ts,
		Payload:    req.Msg.Payload,
	}), nil
}

func (*Echo) ForwardEcho(ctx context.Context, req *connect_go.Request[istioecho.ForwardEchoRequest]) (*connect_go.Response[istioecho.ForwardEchoResponse], error) {
	dest := req.Msg.Url

	res, err := Client(ctx, &EchoClientReq{dest, *req.Msg})

	return connect_go.NewResponse(res), err
}

type EchoClientReq struct {
	Addr string

	Forward istioecho.ForwardEchoRequest
}

func Client(ctx context.Context, req *EchoClientReq) (*istioecho.ForwardEchoResponse, error) {
	return nil, nil
}
