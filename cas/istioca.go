// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cas

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	auth "github.com/costinm/grpc-mesh/gen/grpc-go/istio/v1/auth"
)

const (
	// CertSigner info
	CertSigner = "CertSigner"
)

type Options struct {
	CAEndpoint    string
	CAEndpointSAN string

	TokenProvider credentials.PerRPCCredentials
	GRPCOptions   []grpc.DialOption

	CertSigner string
	ClusterID  string

	CARootPEM    []byte
	TrustedRoots *x509.CertPool

	// ProvCert contains a long-lived 'provider' certificate that will be
	// exchanged with the workload certificate.
	// It is a cert signed by same CA (or a CA trusted by Istiod).
	// It is still exchanged because Istiod may add info to the cert.
	ProvCert string
}

type CitadelClient struct {
	enableTLS bool
	client    auth.IstioCertificateServiceClient
	conn      *grpc.ClientConn
	opts      *Options
}

func trustedCertPool(rootCertPEM []byte) (*x509.CertPool, error) {
	block, rest := pem.Decode(rootCertPEM)
	var blockBytes []byte
	for block != nil {
		blockBytes = append(blockBytes, block.Bytes...)
		block, rest = pem.Decode(rest)
	}

	rootCAs, err := x509.ParseCertificates(blockBytes)
	if err != nil {
		return nil, err
	}
	trustedCertPool := x509.NewCertPool()
	for _, c := range rootCAs {
		trustedCertPool.AddCert(c)
	}
	return trustedCertPool, nil
}

// NewCitadelClient create a CA client for Citadel.
func NewCitadelClient(opts *Options) (*CitadelClient, error) {
	var err error
	if opts.CARootPEM != nil {
		opts.TrustedRoots, err = trustedCertPool(opts.CARootPEM)
		if err != nil {
			return nil, err
		}
	}
	c := &CitadelClient{
		enableTLS: true,
		opts:      opts,
	}

	conn, err := c.buildConnection()
	if err != nil {
		log.Printf("Failed to connect to endpoint %s: %v", opts.CAEndpoint, err)
		return nil, fmt.Errorf("failed to connect to endpoint %s", opts.CAEndpoint)
	}
	c.conn = conn
	c.client = auth.NewIstioCertificateServiceClient(conn)
	return c, nil
}

func (c *CitadelClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// CSR Sign calls Citadel to sign a CSR.
func (c *CitadelClient) CSRSign(csrPEM []byte, certValidTTLInSec int64) ([]string, error) {
	crMetaStruct := &auth.Struct{
		Fields: map[string]*auth.Value{
			CertSigner: {
				Kind: &auth.Value_StringValue{StringValue: c.opts.CertSigner},
			},
		},
	}
	req := &auth.IstioCertificateRequest{
		Csr:              string(csrPEM),
		ValidityDuration: certValidTTLInSec,
		Metadata:         crMetaStruct,
	}
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("ClusterID", c.opts.ClusterID))
	resp, err := c.client.CreateCertificate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create certificate: %v", err)
	}

	if len(resp.CertChain) <= 1 {
		return nil, errors.New("invalid empty CertChain")
	}

	return resp.CertChain, nil
}

func (c *CitadelClient) getTLSDialOption() (grpc.DialOption, error) {
	// Load the TLS root certificate from the specified file.
	// Create a certificate pool
	var certPool *x509.CertPool
	var err error
	if c.opts.TrustedRoots == nil {
		// No explicit certificate - assume the citadel-compatible server uses a public cert
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
	} else {
		certPool = c.opts.TrustedRoots
	}
	var certificate tls.Certificate
	config := tls.Config{
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if c.opts.ProvCert != "" {
				// Load the certificate from disk
				certificate, err = tls.LoadX509KeyPair(
					filepath.Join(c.opts.ProvCert, "cert-chain.pem"),
					filepath.Join(c.opts.ProvCert, "key.pem"))

				if err != nil {
					// we will return an empty cert so that when user sets the Prov cert path
					// but not have such cert in the file path we use the token to provide verification
					// instead of just broken the workflow
					log.Printf("cannot load key pair, using token instead: %v", err)
					return &certificate, nil
				}
				if certificate.Leaf.NotAfter.Before(time.Now()) {
					log.Printf("cannot parse the cert chain, using token instead: %v", err)
					return &tls.Certificate{}, nil
				}
			}
			return &certificate, nil
		},
	}
	config.RootCAs = certPool

	// For debugging on localhost (with port forward)
	// TODO: remove once istiod is stable and we have a way to validate JWTs locally
	if strings.Contains(c.opts.CAEndpoint, "localhost") {
		config.ServerName = "istiod.istio-system.svc"
	}
	if c.opts.CAEndpointSAN != "" {
		config.ServerName = c.opts.CAEndpointSAN
	}

	transportCreds := credentials.NewTLS(&config)
	return grpc.WithTransportCredentials(transportCreds), nil
}

func (c *CitadelClient) buildConnection() (*grpc.ClientConn, error) {
	var ol []grpc.DialOption
	var opts grpc.DialOption
	var err error
	if c.enableTLS {
		opts, err = c.getTLSDialOption()
		if err != nil {
			return nil, err
		}
		ol = append(ol, opts)
	} else {
		opts = grpc.WithInsecure()
		ol = append(ol, opts)
	}
	ol = append(ol, grpc.WithPerRPCCredentials(c.opts.TokenProvider))
	ol = append(ol, c.opts.GRPCOptions...)
	//security.CARetryInterceptor())

	conn, err := grpc.Dial(c.opts.CAEndpoint, ol...)
	if err != nil {
		log.Printf("Failed to connect to endpoint %s: %v", c.opts.CAEndpoint, err)
		return nil, fmt.Errorf("failed to connect to endpoint %s", c.opts.CAEndpoint)
	}

	return conn, nil
}
