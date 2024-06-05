package echo

import (
	"context"
	"errors"
	"log"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect-go/fgrpc/fgrpcconnect"
	"github.com/costinm/grpc-mesh/gen/connect-go/proto/protoconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/fgrpc"
	istioecho "github.com/costinm/grpc-mesh/gen/proto/go/proto"
)

type IstioEcho struct {
	protoconnect.UnimplementedEchoTestServiceHandler
	fgrpcconnect.UnimplementedPingServerHandler
}

func (*IstioEcho) Echo(ctx context.Context, req *connect_go.Request[istioecho.EchoRequest]) (*connect_go.Response[istioecho.EchoResponse], error) {
	log.Println(req.Header(), req.Peer(), req.Spec(), req.Msg)

	res := connect_go.NewResponse(&istioecho.EchoResponse{
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil

	return nil, nil
	//connect_go.NewError(connect_go.CodeUnimplemented, errors.New("proto.EchoTestService.Echo is not implemented"))
}

func (*IstioEcho) Ping(context.Context, *connect_go.Request[fgrpc.PingMessage]) (*connect_go.Response[fgrpc.PingMessage], error) {

	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("fgrpc.PingServer.Ping is not implemented"))
}

func (*IstioEcho) ForwardEcho(ctx context.Context, req *connect_go.Request[istioecho.ForwardEchoRequest]) (*connect_go.Response[istioecho.ForwardEchoResponse], error) {
	dest := req.Msg.Url
	Client(dest, req.Msg)
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("proto.EchoTestService.ForwardEcho is not implemented"))
}

func Client(dest string, msg *istioecho.ForwardEchoRequest) {

}
