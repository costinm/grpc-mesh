package istioca

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect/go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
	"github.com/costinm/meshauth"
	"golang.org/x/net/http2"
)

// Reflective CA.
// Will accept any valid JWT and return a cert for the same identity.
// Trust domain is based on the JWT issuer.
type IstioCA struct {

	// Map of issuers to trust domains.
	TrustDomains map[string]string
	CA           *meshauth.CA
	Auth         *meshauth.MeshAuth

}

func New(ctx context.Context, ma *meshauth.MeshAuth) *IstioCA {


	return &IstioCA{}

}

func (X *IstioCA) CreateCertificate(ctx context.Context, r *connect.Request[auth.IstioCertificateRequest]) (*connect.Response[auth.IstioCertificateResponse], error) {
	slog.Info(":CreateCertificate", "peer", r.Peer(),
		"h", r.Header())

	csrBlock, _ := pem.Decode([]byte(r.Msg.Csr))
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// csrBlock.Type


	slog.Info("csr", "csr", csr)

	a := r.Header().Get("Authorization")
	if a == "" {
		return nil, connect.NewError(http.StatusUnauthorized, errors.New("Missing authorization"))
	}

	j := meshauth.DecodeJWT(a)

	log.Println("CA: ", j)



	tpl := X.CA.CertTemplate("", "")
	crt := X.CA.SignCertificate(tpl, csr.PublicKey.(crypto.PublicKey))

	res := []string{crt, string(X.CA.CACertPEM)}

	return connect.NewResponse(&auth.IstioCertificateResponse{
		CertChain: res,
	}), nil
}


// GetCertIstio implements the basic Istio gRPC protocol as client, getting a cert
// for the primary key in MeshAuth.
//
// The 'dest' must be configured with
//   - TokenSource reading the istio-ca mounted token
//   - K8S token source returning "istio-ca" certs (using cluster, kubeconfig or other user creds)
//   - An existing certificate
//   - A long-lived certificate
func GetCertIstio(ctx context.Context, ma *meshauth.MeshAuth,
	dest *meshauth.Dest, ttlSec int, certSigner string) ([]string, error) {

	// ma.HttpClient returns a client capable of using ma private key.

	//hcPriorK := &http.Client{
	//	// Skip TLS Dial
	//	Transport: &http2.Transport{
	//		AllowHTTP: true,
	//		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
	//			return net.Dial(netw, addr)
	//		},
	//	},
	//}

	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // test server certificate is not trusted.
			},
		},
	}

	caclient := authconnect.NewIstioCertificateServiceClient(// hcPriorK,
		client,
		//http.DefaultClient, // ma.HttpClient(dest),
		"https://" + dest.Addr, connect.WithGRPC())

	csr, err := ma.NewCSR(ma.Cert.PrivateKey, "test.com")

	if ttlSec == 0 {
		ttlSec = 3600
	}

	req := &auth.IstioCertificateRequest{
		Csr:              string(csr),
		ValidityDuration: int64(ttlSec),
	}

	// Istio has an option to specify the 'cert signer' to use.
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

	// Istio also supports setting ClusterID to indicate a JWT signed by another
	// cluster. Parsing issuer should do the same...
	//	creq.Header().Set("ClusterID", kts.MeshCluster.muxID)

	bt, err := dest.TokenGetter(ma).GetToken(ctx, "istio-ca")
	creq.Header().Set("Authorization", "Bearer " + bt)

	certificate, err := caclient.CreateCertificate(context.Background(), creq)
	if err != nil {
		return nil, err
	}
	//log.Println(certificate.Msg)

	return certificate.Msg.CertChain, nil
}

