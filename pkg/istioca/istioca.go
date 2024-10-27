package istioca

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/costinm/grpc-mesh/gen/connect/go/istio/v1/auth/authconnect"
	"github.com/costinm/grpc-mesh/gen/proto/go/istio/v1/auth"
	"github.com/costinm/meshauth"
	meshca "github.com/costinm/meshauth/pkg/ca"
	"github.com/costinm/meshauth/pkg/certs"
	"github.com/costinm/meshauth/pkg/uk8s"
	sshd "github.com/costinm/ssh-mesh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/http2"
)

// IstioCA implements the Istio CA interface as a 'reflective CA'.
//
// Will accept any valid JWT and return a cert for the same identity.
// Trust domain is based on the JWT issuer.
type IstioCA struct {

	// Map of issuers to trust domains.
	IssuerToTrustDomain map[string]string

	CA *meshca.CA

	Auth  *meshauth.Mesh
	SSHCA *sshd.SSHMesh
}

func New(ctx context.Context, ma *meshauth.Mesh) *IstioCA {

	ca := &IstioCA{}
	// Init the SSH CA
	s, err := ssh.NewSignerFromKey(ma.Cert.PrivateKey)
	if err != nil {
		sshCA, _ := sshd.NewSSHMesh(ma)
		// Will be used to sign certificates
		sshCA.Signer = s
		ca.SSHCA = sshCA

	}

	return ca

}

func (ca *IstioCA) SignSSH(public interface{}, id string, secret *uk8s.Secret) {
	sshpk, err := ssh.NewPublicKey(public)
	if err != nil {
		slog.Info("Invalid public", "err", err)
		return
	}
	h, _, err := ca.SSHCA.Sign(sshpk, ssh.HostCert, []string{id})
	secret.Data["ssh-host"] = h
}

// CreateCertificate should be chained with a verification handler.
// The method itself assumes Authorization has been verified.
func (ca *IstioCA) CreateCertificate(ctx context.Context, r *connect.Request[auth.IstioCertificateRequest]) (*connect.Response[auth.IstioCertificateResponse], error) {
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


	if !strings.Contains(r.Msg.Csr, PEMBlockCertificateRequest) {
		// Extended behavior: allow SSH cert signature.

	}

	tpl := ca.CA.CertTemplate("", "")
	crt := ca.CA.SignCertificate(tpl, csr.PublicKey.(crypto.PublicKey))

	res := []string{crt, string(ca.CA.CACertPEM)}

	return connect.NewResponse(&auth.IstioCertificateResponse{
		CertChain: res,
	}), nil
}


// NewCSR creates a key and CSR for the specified SAN (which may be ignored by
// the CA and replaced with a reflexive identity).
// It uses the primary identity in Mesh.
func NewCSR(priv crypto.PrivateKey, san string) (csrPEM []byte, err error) {

	// Exp: Istio overrides the SAN anyways, set it as org to see what happens.
	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{san},
		},
	}

	// Istio sets spiffee - check what happens if a DNS is requested.
	csr.DNSNames = []string{san}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, priv)

	csrPEM = pem.EncodeToMemory(&pem.Block{Type: PEMBlockCertificateRequest, Bytes: csrBytes})

	return
}

const 	PEMBlockCertificateRequest = "CERTIFICATE REQUEST"


// GetCertIstio implements the basic Istio gRPC protocol as client, getting a cert
// for the primary key in Mesh.
//
// The 'dest' must be configured with
//   - TokenSource reading the istio-ca mounted token
//   - K8S token source returning "istio-ca" certs (using cluster, kubeconfig or other user creds)
//   - An existing certificate
//   - A long-lived certificate
func GetCertIstio(ctx context.Context, ma *meshauth.Mesh,
	dest *meshauth.Dest, ttlSec int, certSigner string) (string, []string, error) {

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

	csr, err := NewCSR(ma.Cert.PrivateKey, "test.com")

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
		return "", nil, err
	}
	//log.Println(certificate.Msg)


	return 	string(certs.MarshalPrivateKey(ma.Cert.PrivateKey)), certificate.Msg.CertChain, nil
}

