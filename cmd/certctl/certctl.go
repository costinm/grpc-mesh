package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/gcp"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/mesh"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/sts"
	"github.com/costinm/grpc-mesh/cas"
)

var (
	ns       = flag.String("n", "default", "Namespace")
	aud      = flag.String("audience", "", "Audience to use in the CSR request")
	provider = flag.String("addr", "meshca", "Address. If empty will use the cluster default. meshca or cas can be used as shortcut")
	outDir   = flag.String("out", "/var/run/secrets/workload-spiffe-credentials", "Output dir")
)

// CLI to get the mesh certificates, using CAS or MeshCA.
// The resulting certificate is saved using workload certificate paths.
//
// This should be run periodically to refresh the certs in environments where
// CSI or other native integration is missing.
func main() {
	flag.Parse()

	kr := mesh.New()
	if *ns != "" {
		kr.Namespace = *ns
	}
	ctx := context.Background()

	// Use K8S to get tokens.
	err := gcp.InitGCP(ctx, kr)
	if err != nil {
		log.Fatal("Failed to find K8S ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}

	if false {
		kr.TrustDomain = kr.ProjectId + ".svc.id.goog"
		// Need project number, project_name from cluster
		kr.ProjectNumber = "438684899409"
		kr.ClusterAddress = "https://container.googleapis.com/v1/projects/" +
			kr.ProjectId + "/locations/" + kr.ClusterLocation + "/clusters/" + kr.ClusterName
	} else {
		err = kr.LoadConfig(context.Background())
		if err != nil {
			log.Fatal("Failed to connect to mesh ", time.Since(kr.StartTime), kr, os.Environ(), err)
		}
	}

	// Used to generate the CSR
	auth := sts.NewAuth()
	auth.TrustDomain = kr.TrustDomain
	auth.Namespace = kr.Namespace

	// TODO: fetch public keys too - possibly from all

	cas.InitCA(auth, &sts.AuthConfig{
		ProjectNumber:  kr.ProjectNumber,
		TrustDomain:    kr.TrustDomain,
		ClusterAddress: kr.ClusterAddress,
		TokenSource:    kr.TokenProvider,
	},
		kr.Namespace,
		kr.KSA,
		*provider, "projects/"+kr.ProjectId+
			"/locations/"+kr.Region()+"/caPools/mesh", kr.ClusterLocation,
		kr.MeshConnectorAddr+":15012", kr.CitadelRoot)
}
