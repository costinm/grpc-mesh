package uxds

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect-go/envoy/service/discovery/v3/v3connect"
	"github.com/costinm/grpc-mesh/gen/connect-go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
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

    Id: sidecar~{ip}~{pod_name}.{ns}~{ns}.svc.cluster.local
*/
const WORKLOAD_TYPE = "type.googleapis.com/istio.workload.Workload"

// Resource types in xDS v3.
const (
	apiTypePrefix       = "type.googleapis.com/"
	EndpointType        = apiTypePrefix + "envoy.config.endpoint.v3.ClusterLoadAssignment"
	ClusterType         = apiTypePrefix + "envoy.config.cluster.v3.Dest"
	RouteType           = apiTypePrefix + "envoy.config.route.v3.RouteConfiguration"
	ScopedRouteType     = apiTypePrefix + "envoy.config.route.v3.ScopedRouteConfiguration"
	VirtualHostType     = apiTypePrefix + "envoy.config.route.v3.VirtualHost"
	ListenerType        = apiTypePrefix + "envoy.config.listener.v3.Listener"
	SecretType          = apiTypePrefix + "envoy.extensions.transport_sockets.tls.v3.Secret"
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
			// to delegate to tls.Dial if the network isn't TCP or the addr isn't
			// in an allowlist.
			return net.Dial(network, addr)
		},
		// Don't forget timeouts!
	}
}

func NewXDS(cfg *XDSConfig) *XDS {

	xdsdest := &meshauth.Dest{BaseAddr: cfg.XDS}

	client := v3connect.NewAggregatedDiscoveryServiceClient(xdsdest.H2Client(), xdsdest.BaseAddr, connect.WithGRPC())

	return &XDS{
		client: client,
		cfg:    cfg,
		//		ds: ds,
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
			fmt.Println(
				ws)
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

// GetCertIstio implements the basic Istio gRPC protocol
// The 'dest' may be configured with
//   - TokenSource reading the istio-ca mounted token
//   - K8S token source returning "istio-ca" certs (using cluster, kubeconfig or other user creds)
//   - An existing certificate
//   - A long-lived certificate
func GetCertIstio(ctx context.Context, dest *meshauth.Dest, ttlSec int, certSigner string) ([]byte, []string, error) {
	caclient := authconnect.NewIstioCertificateServiceClient(dest.H2Client(),
		dest.BaseAddr, connect.WithGRPC())

	priv, csr, err := dest.MeshAuth.NewCSR("")

	req := &auth.IstioCertificateRequest{
		Csr:              string(csr),
		ValidityDuration: int64(ttlSec),
	}
	if certSigner != "" {
		req.Metadata = &auth.Struct{
			Fields: map[string]*auth.Value{
				"CertSigner": {
					Kind: &auth.Value_StringValue{StringValue: certSigner},
				},
			},
		}
	}

	creq := connect.NewRequest(req)

	//if kts, ok := dest.TokenSource.(*meshauth.K8sTokenSource); ok {
	//	// This is complicated: most of the type default cluster is named 'Kubernetes',
	//	// and in the rare 'one istiod handling multiple clusters' the names on Istiod
	//	// side are arbitrary.
	//	creq.Header().Set("ClusterID", kts.Dest.ID)
	//}
	bt, err := dest.TokenSource.GetRequestMetadata(ctx, "istio-ca")
	for k, v := range bt {
		creq.Header().Set(k, v)
	}

	certificate, err := caclient.CreateCertificate(context.Background(), creq)
	if err != nil {
		return nil, nil, err
	}
	log.Println(certificate.Msg)

	return priv, certificate.Msg.CertChain, nil
}
