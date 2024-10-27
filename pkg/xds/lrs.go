
package xds

import (
	"context"

	connect_go "github.com/bufbuild/connect-go"
	lrspb "github.com/costinm/grpc-mesh/gen/proto/go/envoy/service/load_stats/v2"
	"github.com/golang/protobuf/ptypes/duration"

	"log"
	"sync"
)

type NodeInfo struct {
	ClusterStats []*lrspb.ClusterStats
}

type LRS struct {
	m sync.Mutex

	nodeInfo map[string]*NodeInfo
}

func NewLRS() *LRS {
	return &LRS{}
}

func (l *LRS) Dump() {

}

func (l *LRS) StreamLoadStats(ctx context.Context, s *connect_go.BidiStream[lrspb.LoadStatsRequest, lrspb.LoadStatsResponse]) error {
	log.Println(s.Peer())
	log.Println(s.RequestHeader())

	// Send all clusters - instead of selective list
	s.Send(&lrspb.LoadStatsResponse{SendAllClusters: true,
		LoadReportingInterval:     &duration.Duration{Seconds: 60},
		ReportEndpointGranularity: true,
	})
	for {
		req, err := s.Receive()
		if err != nil {
			return err
		}
		log.Print(req.Node.Id)
		for _, lr := range req.ClusterStats {
			log.Println(lr.ClusterName)
		}
		l.m.Lock()
		l.nodeInfo[req.Node.Id] = &NodeInfo{ClusterStats: req.ClusterStats}
	}

	return nil
}
