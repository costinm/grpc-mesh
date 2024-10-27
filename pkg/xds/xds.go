package xds

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect-go/envoy/service/discovery/v3/v3connect"
	"github.com/costinm/grpc-mesh/gen/proto/go/xds"
	"github.com/costinm/meshauth"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

/*
Required node info in ambient:

  - INSTANCE_IPS ( env var is INSTANCE_IP )

  - POD_NAME

  - NAMESPACE - from POD_NAMESPACE

  - NODE_NAME

  - AMBIENT_TYPE

    muxID: sidecar~{ip}~{pod_name}.{ns}~{ns}.svc.cluster.local
*/
const WORKLOAD_TYPE = "type.googleapis.com/istio.workload.Workload"

// istio/pkg/model
const (
	APITypePrefix   = "type.googleapis.com/"
	envoyTypePrefix = APITypePrefix + "envoy."

	ClusterType                = APITypePrefix + "envoy.config.cluster.v3.Cluster"
	EndpointType               = APITypePrefix + "envoy.config.endpoint.v3.ClusterLoadAssignment"
	ListenerType               = APITypePrefix + "envoy.config.listener.v3.Listener"
	RouteType                  = APITypePrefix + "envoy.config.route.v3.RouteConfiguration"
	SecretType                 = APITypePrefix + "envoy.extensions.transport_sockets.tls.v3.Secret"
	ExtensionConfigurationType = APITypePrefix + "envoy.config.core.v3.TypedExtensionConfig"

	NameTableType   = APITypePrefix + "istio.networking.nds.v1.NameTable"
	HealthInfoType  = APITypePrefix + "istio.v1.HealthInformation"
	ProxyConfigType = APITypePrefix + "istio.mesh.v1alpha1.ProxyConfig"
	// DebugType requests debug info from istio, a secured implementation for istio debug interface.
	DebugType                 = "istio.io/debug"
	BootstrapType             = APITypePrefix + "envoy.config.bootstrap.v3.Bootstrap"
	AddressType               = APITypePrefix + "istio.workload.Address"
	WorkloadType              = APITypePrefix + "istio.workload.Workload"
	WorkloadAuthorizationType = APITypePrefix + "istio.security.Authorization"
)

// Resource types in xDS v3.
const (
	apiTypePrefix       = "type.googleapis.com/"
	ScopedRouteType     = apiTypePrefix + "envoy.config.route.v3.ScopedRouteConfiguration"
	VirtualHostType     = apiTypePrefix + "envoy.config.route.v3.VirtualHost"
	ExtensionConfigType = apiTypePrefix + "envoy.config.core.v3.TypedExtensionConfig"
	RuntimeType         = apiTypePrefix + "envoy.service.runtime.v3.Runtime"

	// AnyType is used only by ADS
	AnyType = ""
)

const RBAC_TYPE = "type.googleapis.com/istio.workload.Authorization"

type XDSConfig struct {
	// Namespace defaults to 'default'
	Namespace string

	// Workload defaults to 'test'
	Workload string

	XDSHeaders map[string]string

	// IP is currently the primary key used to locate inbound configs. It is sent by client,
	// must match a known endpoint IP. Tests can use a ServiceEntry to register fake IPs.
	IP string

	// Context used for early cancellation
	Context context.Context

	// Base URL of the XDS server, including scheme
	XDS string
	// Not included: Locality, Meta, GrpcOpts, NodeId

	MeshAuth *meshauth.Mesh
}

type XDS struct {
	ds     StreamService[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse] //*connect.BidiStreamForClient[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse]
	fs     StreamService[xds.DiscoveryRequest, xds.DiscoveryResponse]           //*connect.BidiStreamForClient[xds.DeltaDiscoveryRequest, xds.DeltaDiscoveryResponse]
	client v3connect.AggregatedDiscoveryServiceClient
	cfg    *XDSConfig
}

type StreamService[I any, O any] interface {
	Receive() (*O, error)
	Send(*I) error
}

var (
	XDSTopics map[string]Topic = map[string]Topic{
		"ptr": &XDSTopic[*xds.Workload]{
			TypeURL: WORKLOAD_TYPE,
			T:       &xds.Workload{},
		},
		"cluster": &XDSTopic[*xds.Cluster]{
			TypeURL: ClusterType,
			T:       &xds.Cluster{},
		},
	}
)

type Topic interface {
	Proto() proto.Message
	GetTypeURL() string
}

type XDSTopic[T proto.Message] struct {
	TypeURL       string
	T             T
	ResourceNames map[string]string
	Resourcces    map[string]T
}

func (x *XDSTopic[T]) Proto() proto.Message {
	return x.T
}
func (x *XDSTopic[T]) GetTypeURL() string {
	return x.TypeURL
}

type XDSResource[T any] struct {
	Value   T
	Name    string
	TypeURL string
}

func TransportFunc(d *meshauth.Dest) http.RoundTripper {
	return &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			// If you're also using this client for non-h2c traffic, you may want
			// to delegate to tls.DialContext if the network isn't TCP or the addr isn't
			// in an allowlist.
			return net.Dial(network, addr)
		},
		// Don't forget timeouts!
	}
}

// NewXDS creates a new XDS client
func NewXDS(cfg *XDSConfig) *XDS {

	xdsdest := &meshauth.Dest{Addr: cfg.XDS}
	url := "https://" + cfg.XDS
	if strings.HasSuffix(cfg.XDS, ":15010") {
		xdsdest.L4Secure = true
		url = "http://" + cfg.XDS
	}
	hc2 := &http.Client{
		// Skip TLS Dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	client := v3connect.NewAggregatedDiscoveryServiceClient(//cfg.Mesh,
		hc2, // HttpClient(xdsdest),
		url, connect.WithGRPC())

	return &XDS{
		client: client,
		cfg:    cfg,
	}
}

func (x *XDS) RunFull(initial string) error {
	cfg := x.cfg

	t := XDSTopics[initial]
	if t == nil {
		return nil
	}
	req := &xds.DiscoveryRequest{
		TypeUrl: t.GetTypeURL(),
		Node: &xds.Node{Id: fmt.Sprintf("sidecar~%s~%s.%s~%s.svc.cluster.local",
			cfg.IP, cfg.Workload, cfg.Namespace, cfg.Namespace)},
	}
	ctx := context.Background()

	ds := x.client.StreamAggregatedResources(ctx)
	ds.RequestHeader().Set("Test", "header")
	// HTTP request on first Send.
	err := ds.Send(req)
	if err != nil {
		return err
	}

	for _, s := range []string{} {
		req = &xds.DiscoveryRequest{
			TypeUrl: apiTypePrefix + s,
		}
		err = ds.Send(req)
		if err != nil {
			return err
		}

	}

	x.fs = ds
	for {
		res, err := x.fs.Receive()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("type", res.TypeUrl, "nonce", res.Nonce)

		x.fs.Send(&xds.DiscoveryRequest{
			TypeUrl:       res.TypeUrl,
			ResponseNonce: res.Nonce})

		for _, r := range res.Resources {
			w := t.Proto().ProtoReflect().New().Interface()
			proto.Unmarshal(r.Value, w)
			ws := prototext.MarshalOptions{Multiline: false}.Format(w)
			fmt.Println("res", ws)
		}
	}
}

func (x *XDS) RunDelta(initial string) error {
	cfg := x.cfg

	t := XDSTopics[initial]
	if t == nil {
		return nil
	}

	req := &xds.DeltaDiscoveryRequest{
		TypeUrl: t.GetTypeURL(),
		Node: &xds.Node{Id: fmt.Sprintf("sidecar~%s~%s.%s~%s.svc.cluster.local",
			cfg.IP, cfg.Workload, cfg.Namespace, cfg.Namespace)},
	}
	ctx := context.Background()

	ds := x.client.DeltaAggregatedResources(ctx)
	err := ds.Send(req)
	if err != nil {
		return err
	}
	x.ds = ds
	for {
		res, err := x.ds.Receive()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("type", res.TypeUrl, "nonce", res.Nonce,
			"systemVersionInfo", res.SystemVersionInfo, "removed", res.RemovedResources)

		x.ds.Send(&xds.DeltaDiscoveryRequest{
			TypeUrl:       res.TypeUrl,
			ResponseNonce: res.Nonce})

		for _, r := range res.Resources {
			w := t.Proto().ProtoReflect().New().Interface()
			proto.Unmarshal(r.Resource.Value, w)
			ws := prototext.MarshalOptions{Multiline: false}.Format(w)
			fmt.Println("name:", r.Name,
				//"aliases", r.Aliases,
				//"ttl", r.Ttl,
				//"version", r.Version,
				//"cache-control", r.CacheControl,
				//"typeUrl", r.Resource.TypeUrl,
				ws)
		}
	}
}

