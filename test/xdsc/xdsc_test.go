package xdsc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/k8s"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/mesh"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/sts"
	"github.com/costinm/grpc-mesh/xdsc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// Required for k8s client to link in the authenticator
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Get mesh config from the current cluster
func LoadConfig(ctx context.Context) (*mesh.KRun, error) {
	kr := mesh.New()

	// Using a K8s Client
	kc := &k8s.K8S{Mesh: kr}
	err := kc.K8SClient(ctx)
	if err != nil {
		return nil, err
	}
	kr.Cfg = kc
	kr.TokenProvider = kc

	// Load the config map with mesh-env
	kr.LoadConfig(ctx)

	return kr, err
}

func TestMCP(t *testing.T) {

	ctx, cf := context.WithTimeout(context.Background(), 100*time.Second)
	defer cf()

	kr, err := LoadConfig(ctx)

	// The mesh configuration should have all the properties we need -
	// see the code for example content.

	m := map[string]interface{}{}
	m["CLUSTER_ID"] = kr.ClusterName
	m["NAMESPACE"] = kr.Namespace
	m["SERVICE_ACCOUNT"] = kr.KSA
	// Required for connecting to MCP.
	m["CLOUDRUN_ADDR"] = "asm-big1-asm-managed-us-central1-c-42okyzkgcq-uc.a.run.app:443"
	// kr.MeshTenant

	//m["LABELS"] = x.Labels

	m["ISTIO_VERSION"] = "1.20.0-xdsc"
	m["SDS"] = "true"

	xdsConfig := &xdsc.Config{
		Namespace: kr.Namespace,
		Workload:  kr.Name + "-" + "10.10.1.1",
		Meta:      m,
		NodeType:  "sidecar",
		IP:        "10.10.1.1",
		Context:   ctx,

		GrpcOpts: []grpc.DialOption{
			// Using the STS library to exchange tokens
			grpc.WithPerRPCCredentials(sts.NewGSATokenSource(kr, "")),
			grpc.WithTransportCredentials(credentials.NewTLS(
				&tls.Config{
					InsecureSkipVerify: false,
				})),
		},
	}

	//ctx = metadata.AppendToOutgoingContext(context.Background(), "ClusterID", p.clusterID)
	xdsAddr := kr.XDSAddr
	xdsAddr = "staging-meshconfig.sandbox.googleapis.com:443"
	// "meshconfig.googleapis.com:443"
	xdscc, err := xdsc.DialContext(ctx, xdsAddr, xdsConfig)
	// calls Run()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Connected", xdscc)
	xdscc.Watch()

	xdscc.Fetch()
	log.Println(xdscc.Responses)
}

func TestIstiod(t *testing.T) {

	ctx, cf := context.WithTimeout(context.Background(), 100*time.Second)
	defer cf()

	kr, err := LoadConfig(ctx)

	// The mesh configuration should have all the properties we need -
	// see the code for example content.

	m := map[string]interface{}{}
	m["CLUSTER_ID"] = kr.ClusterName
	m["NAMESPACE"] = kr.Namespace
	m["SERVICE_ACCOUNT"] = kr.KSA
	// kr.MeshTenant

	//m["LABELS"] = x.Labels

	m["ISTIO_VERSION"] = "1.20.0-xdsc"
	m["SDS"] = "true"

	xdsConfig := &xdsc.Config{
		Namespace: kr.Namespace,
		Workload:  kr.Name + "-" + "10.10.1.1",
		Meta:      m,
		NodeType:  "sidecar",
		IP:        "10.10.1.1",
		Context:   ctx,

		GrpcOpts: []grpc.DialOption{
			// Using the STS library to exchange tokens
			grpc.WithPerRPCCredentials(
				sts.NewK8STokenSource(kr, kr.TrustDomain)),
			grpc.WithTransportCredentials(credentials.NewTLS(
				&tls.Config{
					// TODO: specify the cert from kr
					InsecureSkipVerify: true,
				})),
		},
	}

	//ctx = metadata.AppendToOutgoingContext(context.Background(), "ClusterID", p.clusterID)
	xdsAddr := kr.XDSAddr
	if xdsAddr == "" {
		// This is the external address, only enabled for multi-network.
		// If running from an internal address, use
		// kr.MeshConnectorInternalAddr
		xdsAddr = kr.MeshConnectorAddr + ":15012"
	}

	xdscc, err := xdsc.DialContext(ctx, xdsAddr, xdsConfig)
	// calls Run()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Connected", xdscc)
	xdscc.Watch()

	xdscc.Fetch()
	log.Println(xdscc.Responses)
}

func TestTD(t *testing.T) {

	ctx, cf := context.WithTimeout(context.Background(), 100*time.Second)
	defer cf()

	kr, err := LoadConfig(ctx)

	m := map[string]interface{}{}
	m["TRAFFICDIRECTOR_NETWORK_NAME"] = "default"
	m["TRAFFICDIRECTOR_GCP_PROJECT_NUMBER"] = kr.ProjectNumber
	m["INSTANCE_IP"] = "10.48.0.63"

	xdsConfig := &xdsc.Config{
		Namespace: kr.Namespace,
		Workload:  kr.Name + "-" + "10.10.1.1",
		Meta:      m,

		NodeType: "sidecar",
		IP:       "10.10.1.1",
		Context:  ctx,
		Locality: kr.ClusterLocation,
		NodeId: fmt.Sprintf("projects/%s/networks/default/nodes/1234",
			kr.ProjectNumber),

		GrpcOpts: []grpc.DialOption{
			// Using the STS library to exchange tokens
			grpc.WithPerRPCCredentials(
				sts.NewGSATokenSource(kr, "k8s-fortio@costin-asm1.iam.gserviceaccount.com")),
			grpc.WithTransportCredentials(credentials.NewTLS(
				&tls.Config{
					// TODO: specify the cert from kr
					InsecureSkipVerify: true,
				})),
		},
	}

	xdscc, err := xdsc.DialContext(ctx, "trafficdirector.googleapis.com:443", xdsConfig)
	// calls Run()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Connected", xdscc)
	xdscc.Watch()

	xdscc.Fetch()
	log.Println(xdscc.Responses)
}
