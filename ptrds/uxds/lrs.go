package uxds

import (
	"context"
	"errors"
	connect_go "github.com/bufbuild/connect-go"
	lrspb "github.com/costinm/grpc-mesh/gen/proto/go/envoy/service/load_stats/v2"
	"log"
)

type LRS struct {
}

func NewLRS() *LRS {
	return &LRS{}
}

func (l *LRS) StreamLoadStats(ctx context.Context, s *connect_go.BidiStream[lrspb.LoadStatsRequest, lrspb.LoadStatsResponse]) error {
	log.Println(s.Peer())
	log.Println(s.RequestHeader())
	for {
		req, err := s.Receive()
		if err != nil {
			return err
		}
		log.Print(req.Node.Id)
		for _, lr := range req.ClusterStats {
			log.Println(lr.ClusterName)
		}
	}

	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("envoy.service.load_stats.v2.LoadReportingService.StreamLoadStats is not implemented"))
}
