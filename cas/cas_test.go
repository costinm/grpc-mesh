package cas_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/costinm/grpc-mesh/bootstrap"
	"github.com/costinm/grpc-mesh/cas"
	"github.com/costinm/meshauth"
	"google.golang.org/grpc"
)

func TestCAs(t *testing.T) {
	kconf, err := bootstrap.LoadKubeconfig()
	if err != nil {
		t.Skip("Can't find a kube config file")
	}

	ctx := context.Background()
	def, _, err := meshauth.InitK8S(ctx, kconf)
	if err != nil {
		t.Fatal(err)
	}

	ma := meshauth.NewMeshAuth(nil)
	priv, csr, err := ma.NewCSR("spiffe://cluster.local/ns/" + def.Namespace + "/sa/" +
		def.ServiceAccount)
	if err != nil {
		t.Fatal("Failed to find mesh certificates ", err)
	}

	t.Run("Istio", func(t *testing.T) {
		// Tokens using istio-ca audience for Istio - this is what Citadel and Istiod expect
		catokenS := def.NewK8STokenSource("istio-ca")
		//catokenS := def.NewK8STokenSource(def.ProjectID + ".svc.id.goog")

		caroot := def.GetEnv("CAROOT_ISTIOD", "")
		addr := def.GetEnv("MCON_ADDR", "")
		if addr == "" {
			addr = "127.0.0.1:15012"
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
			t.Fatal(err)
		}
		chain, err := cca.CSRSign(csr, 24*3600)
		if err != nil {
			t.Fatal(err)
		}
		chainPEMCat := strings.Join(chain, "\n")
		tlsCert, err := tls.X509KeyPair([]byte(chainPEMCat), priv)
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range tlsCert.Certificate {
			certF, err := x509.ParseCertificate(c)
			if err != nil {
				panic(err)
			}

			log.Println("uri", certF.URIs, "sub", certF.Subject, "iss", certF.Issuer, "issurl", certF.IssuingCertificateURL)
		}
	})

	t.Run("MeshCA", func(t *testing.T) {
		sts1, err := def.GCPFederatedSource(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var ol2 []grpc.DialOption
		ol2 = append(ol2, grpc.WithPerRPCCredentials(sts1))

		mca, err := cas.NewGoogleCAClient("meshca.googleapis.com:443", ol2)
		// Required
		mca.Location = def.Location

		chain, err := mca.CSRSign(csr, 24*3600)
		if err != nil {
			t.Fatal(err)
		}

		chainPEMCat := strings.Join(chain, "\n")
		tlsCert, err := tls.X509KeyPair([]byte(chainPEMCat), priv)
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range tlsCert.Certificate {
			certF, err := x509.ParseCertificate(c)
			if err != nil {
				panic(err)
			}

			log.Println("uri", certF.URIs, "sub", certF.Subject, "iss", certF.Issuer, "issurl", certF.IssuingCertificateURL)
		}
	})

	t.Run("CAS", func(t *testing.T) {
		sts1, err := def.GCPFederatedSource(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var ol2 []grpc.DialOption
		ol2 = append(ol2, grpc.WithPerRPCCredentials(sts1))

		// TODO: use env variable to override, to avoid one RTT to K8S
		// There is no need to configure this in cluster for API-based cert init
		res, err := def.Do(def.Request(ctx, "",
			"security.cloud.google.com/v1/workloadcertificateconfigs", "default", nil))
		if err != nil {
			t.Fatal(err)
		}
		wc := cas.WorkloadCertificateConfig{}
		json.Unmarshal(res, &wc)
		log.Println(wc)

		// This must use region, not location
		//mca, err := NewGoogleCASClientRaw("projects/costin-asm1/locations/us-central1/caPools/mesh",
		//	ol2)

		pool := strings.TrimPrefix(wc.Spec.CertificateAuthorityConfig.CertificateAuthorityServiceConfig.EndpointURI,
			"//privateca.googleapis.com/")

		pool1 := fmt.Sprintf("projects/%s/locations/%s/caPools/mesh", def.ProjectID,
			cas.Region(def.Location))
		if pool1 != pool {
			log.Println("Invalid pool", pool1, pool)
		}
		mca, err := cas.NewGoogleCASClientRaw(pool1,
			ol2)

		chain, err := mca.CSRSign(csr, 24*3600)
		if err != nil {
			t.Fatal(err)
		}

		chainPEMCat := strings.Join(chain, "\n")
		tlsCert, err := tls.X509KeyPair([]byte(chainPEMCat), priv)
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range tlsCert.Certificate {
			certF, err := x509.ParseCertificate(c)
			if err != nil {
				panic(err)
			}

			log.Println("uri", certF.URIs, "sub", certF.Subject, "iss", certF.Issuer, "issurl", certF.IssuingCertificateURL)
		}
	})
}
