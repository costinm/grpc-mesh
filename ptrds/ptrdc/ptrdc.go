package main

import (
	"context"
	"flag"
	"log"

	"github.com/costinm/grpc-mesh/ptrds/uxds"
)

var (
	IstiodAddr = flag.String("xds", "http://localhost:15010", "Istio address")
	ns         = flag.String("n", "istio-system", "")
	ip         = flag.String("ip", "127.1.2.3", "")
	pod        = flag.String("pod", "test-abc-def", "")
)

// Client for ZTunnel PTR-DS, using connect-go implementation for evaluation.
func main() {
	x := uxds.NewXDS(&uxds.XDSConfig{
		Namespace:  *ns,
		Workload:   *pod,
		XDSHeaders: nil,
		IP:         *ip,
		XDS:        *IstiodAddr,
		Context:    context.Background(),
	})

	//go func() {
	//	err := x.RunDelta("ptr")
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}()

	//err := x.RunFull("cluster")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//err := x.RunDelta("cluster")
	err := x.RunDelta("ptr")
	if err != nil {
		log.Fatalln(err)
	}
}
