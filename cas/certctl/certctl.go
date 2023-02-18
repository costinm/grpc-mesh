package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/costinm/grpc-mesh/bootstrap"
	"github.com/costinm/grpc-mesh/cas"
	"github.com/costinm/meshauth"
)

var (
	ns  = flag.String("n", "default", "Namespace")
	aud = flag.String("audience", "", "Audience to use in the CSR request")

	caAddr = flag.String("addr", os.Getenv("CA_ADDR"), "Address. If empty will use the cluster default. meshca "+
		"or cas can be used as shortcut")

	outDir = flag.String("out", "/var/run/secrets/workload-spiffe-credentials", "Output dir")

	casPool = flag.String("cas", "", "If set, use the CAS pool. Can be full name (projects/../locations/../caPools/NAME or just name. Usually 'mesh' name should be used.")
)

// CLI to get the mesh certificates, using Istio, CAS or MeshCA.
//
// The resulting certificate is saved using workload certificate paths.
//
// This should be run periodically to refresh the certs in environments where
// CSI or other native integration is missing.
//
// As source of trust:
//   - in-cluster K8S or KUBECONFIG. For in-cluster, either permissions to call TokenRequest or
//     a mounted token with istio-ca or the trust domain as audience
//   - a GCP metadata or service account file
func main() {
	flag.Parse()

	ctx := context.Background()

	// Start from kubeconfig first
	kconf, err := bootstrap.LoadKubeconfig()

	def, _, err := meshauth.InitK8S(ctx, kconf)
	if err != nil {
		log.Fatal(err)
	}

	// Used to generate the CSR
	auth := meshauth.NewMeshAuth(&meshauth.MeshAuthCfg{
		Namespace: *ns,
	})
	//	auth.TrustDomain = kr.TrustDomain
	//	auth.Namespace = kr.Namespace

	// TODO: fetch public keys too - possibly from all
	priv, csr, err := auth.NewCSR("spiffe://cluster.local/ns/" + def.Namespace + "/sa/" +
		def.ServiceAccount)
	if err != nil {
		log.Fatal("Failed to find mesh certificates ", err)
	}

	// If CAS_POOL is set - use CAS
	// If ProjectNumber is available, use MeshCA
	// else - istio

	var chain []string
	if *caAddr == "cas" || true {
		sts1, err := def.GCPFederatedSource(ctx)
		if err != nil {
			log.Fatal(err)
		}
		pool1 := *casPool
		if *casPool == "" {
			pool1 = "mesh"
		}
		if !strings.Contains(pool1, "/") {
			pool1 = fmt.Sprintf("projects/%s/locations/%s/caPools/%s", def.ProjectID,
				cas.Region(def.Location), pool1)
		}
		mca, err := cas.NewGoogleCASClientRaw(pool1, sts1)
		if err != nil {
			log.Fatal(err)
		}
		chain, err = mca.CSRSign(csr, 24*3600)
		if err != nil {
			log.Fatal(err)
		}

	} else if *caAddr == "meshca" {
		sts1, err := def.GCPFederatedSource(ctx)
		if err != nil {
			log.Fatal(err)
		}
		mca, err := cas.NewGoogleCAClient("meshca.googleapis.com:443", sts1)
		// Required
		mca.Location = def.Location

		chain, err = mca.CSRSign(csr, 24*3600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		catokenS := def.NewK8STokenSource("istio-ca")

		caroot := def.GetEnv("CAROOT_ISTIOD", "")

		addr := def.GetEnv("MCON_ADDR", "")
		if addr == "" {
			addr = "istiod.istio-system.svc:15012"
		} else {
			addr = addr + ":15012"
		}
		cca, err := cas.NewCitadelClient(&cas.Options{
			TokenProvider: catokenS,
			CAEndpoint:    addr,
			CARootPEM:     []byte(caroot),
			CAEndpointSAN: "istiod.istio-system.svc",
		})
		if err != nil {
			log.Fatal(err)
		}

		chain, err = cca.CSRSign(csr, 24*3600)
		if err != nil {
			log.Fatal(err)
		}
	}

	chainPEMCat := strings.Join(chain, "\n")
	tlsCert, err := tls.X509KeyPair([]byte(chainPEMCat), priv)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range tlsCert.Certificate {
		certF, err := x509.ParseCertificate(c)
		if err != nil {
			panic(err)
		}

		log.Println("uri", certF.URIs, "sub", certF.Subject, "iss", certF.Issuer, "issurl", certF.IssuingCertificateURL)
	}

	auth.SetTLSCertificate(&tlsCert)

	err = auth.SaveCerts(".")
	if err != nil {
		log.Fatal(err)
	}

}
