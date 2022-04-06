//// Copyright Istio Authors
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////     http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
//
package xdsc

//
//import (
//	"context"
//	"flag"
//	"log"
//	"strings"
//
//	"github.com/costinm/istio-grpc/bootstrap"
//	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/balancer"
//	"google.golang.org/grpc/credentials/tls/certprovider"
//	"google.golang.org/grpc/grpclog"
//	"google.golang.org/grpc/resolver"
//	"google.golang.org/grpc/serviceconfig"
//
//	//  To install the xds resolvers and balancers.
//	_ "google.golang.org/grpc/xds"
//)
//
//var (
//	grpcAddr = "127.0.0.1:14057"
//
//	// Address of the Istiod gRPC service, used in tests.
//	istiodSvcAddr = "istiod.istio-system.svc.cluster.local:15010"
//)
//
//type lv2 struct {
//
//}
//
//func (l lv2) Info(args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Infoln(args ...interface{}) {
//	if len(args) < 2 {
//		return
//	}
//	a0 := args[0].(string)
//	if a0 != "[xds]" {
//		return
//	}
//	a1 := args[1].(string)
//	if strings.HasPrefix(a1, "[xds-client") {
//		log.Println(args...)
//	} else if strings.HasPrefix(a1, "[eds-lb") {
//		log.Println(args...)
//	} else {
//		//log.Println(args...)
//	}
//}
//
//func (l lv2) Infof(format string, args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Warning(args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Warningln(args ...interface{}) {
//	log.Println(args...)
//}
//
//func (l lv2) Warningf(format string, args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Error(args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Errorln(args ...interface{}) {
//	log.Println(args...)
//}
//
//func (l lv2) Errorf(format string, args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Fatal(args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Fatalln(args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) Fatalf(format string, args ...interface{}) {
//	panic("implement me")
//}
//
//func (l lv2) V(lv int) bool {
//	if lv != 2 {
//		log.Output(2, "level")
//	}
//	return true
//}
//var (
//	addr = flag.String("addr", "istiod-v110.istio-system.svc.cluster.local:15010", "Address to dial")
//)
//
//func main() {
//	// Can't be set - grpc is loading it at init time. Should have some default...
//	//os.Setenv("GRPC_XDS_BOOTSTRAP", "xds_bootstrap.json")
//	flag.Parse()
//
//	err := bootstrap.GenerateBootstrap()
//
//	if err != nil {
//		log.Panic("GRPC_XDS_BOOTSTRAP must be set")
//	}
//
//	grpclog.SetLoggerV2(&lv2{})
//
//	resolve(*addr)
//
//	//resolve(istiodSvcAddr)
//	dialXDS("xds:///" + *addr)
//}
//
//
//func resolve(name string) {
//	rb := resolver.Get("xds")
//		stateCh := &Channel{ch: make(chan interface{}, 1)}
//		errorCh := &Channel{ch: make(chan interface{}, 1)}
//
//		// Depends on version -
//		res, err := rb.Build(resolver.Target{
//			Endpoint: name},
//			&testClientConn{stateCh: stateCh, errorCh: errorCh},
//			resolver.BuildOptions{
//
//			})
//
//		if err != nil {
//			log.Fatal("Failed to resolve XDS ", err)
//		}
//		res.ResolveNow(resolver.ResolveNowOptions{})
//
//		//tm := time.After(5 * time.Second)
//		go func() {
//			for {
//				select {
//				case s := <-stateCh.ch:
//					rs := s.(resolver.State)
//					log.Println("Resolved ", rs.Addresses, rs.ServiceConfig, rs.Attributes)
//					if len(rs.Addresses) > 0 {
//						return
//					}
//					// rs.Attributes has csKeyType with configSelector value, of type xdsResolver
//
//				case e := <-errorCh.ch:
//					log.Println("Error in resolve", e)
//					return
//				//case <-tm:
//				//	log.Println("Didn't resolve")
//				//	return
//				}
//			}
//		}()
//}
//
//// Example: grpc/credentials/tls/certprovider/pemfile
//// Using name=file_watcher
//// {
////				"certificate_file":   "/a/b/cert.pem",
////				"private_key_file":    "/a/b/key.pem",
////				"ca_certificate_file": "/a/b/ca.pem",
////				"refresh_interval":   "200s"
////			}`),
////
////			wantOutput: "file_watcher:/a/b/cert.pem:/a/b/key.pem:/a/b/ca.pem:3m20s"
//type IstioCertProvider struct {
//
//}
//
//func (i IstioCertProvider) ParseConfig(i2 interface{}) (*certprovider.BuildableConfig, error) {
//	panic("implement me")
//}
//
//func (i IstioCertProvider) Name() string {
//	panic("implement me")
//}
//
//func certProv() {
//	certprovider.Register(&IstioCertProvider{})
//}
//
//func balance() {
//		rb := balancer.Get("eds_experimental")
//		b := rb.Build(&testLBClientConn{}, balancer.BuildOptions{})
//		defer b.Close()
//}
//
//// connect to the XDS server, using gRPC
//func dialXDS(addr string) {
//		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//		//defer cancel()
//		ctx := context.Background()
//
//		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
//		if err != nil {
//			log.Fatal("XDS gRPC", err)
//		}
//
//		defer conn.Close()
//		xds := discovery.NewAggregatedDiscoveryServiceClient(conn)
//
//		s, err := xds.StreamAggregatedResources(ctx)
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.Println(s.Send(&discovery.DiscoveryRequest{}))
//
//	for {
//		msg, err := s.Recv()
//		if err != nil {
//			log.Println("XDS error", err)
//			return
//		}
//		log.Println("XDS msg", msg)
//	}
//
//	}
//
//type testLBClientConn struct {
//	balancer.ClientConn
//}
//
//type Channel struct {
//	ch chan interface{}
//}
//
//// Send sends value on the underlying channel.
//func (c *Channel) Send(value interface{}) {
//	c.ch <- value
//}
//
//// From xds_resolver_test
//// testClientConn is a fake implemetation of resolver.ClientConn. All is does
//// is to store the state received from the resolver locally and signal that
//// event through a channel.
//type testClientConn struct {
//	resolver.ClientConn
//	stateCh *Channel
//	errorCh *Channel
//}
//
//func (t *testClientConn) UpdateState(s resolver.State) error{
//	t.stateCh.Send(s)
//	return nil
//}
//
//func (t *testClientConn) ReportError(err error) {
//	t.errorCh.Send(err)
//}
//
//func (t *testClientConn) ParseServiceConfig(jsonSC string) *serviceconfig.ParseResult {
//	// Will be called with something like:
//	//
//	//	"loadBalancingConfig":[
//	//	{
//	//		"cds_experimental":{
//	//			"Cluster": "istiod.istio-system.svc.cluster.local:14056"
//	//		}
//	//	}
//	//]
//	//}
//	return &serviceconfig.ParseResult{}
//}
