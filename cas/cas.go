// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cas

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	// Just protos, no extra deps (grpc, protobuf)
	//privatecapb "google.golang.org/genproto/googleapis/cloud/security/privateca/v1"
	privatecapb "cloud.google.com/go/security/privateca/apiv1/privatecapb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

// GoogleCASClient: Agent or user space plugin for Google CAS
//

type GoogleCASClient struct {
	capool      string
	caClientRaw privatecapb.CertificateAuthorityServiceClient
}

// CAS is defined as grpc service in https://github.com/googleapis/go-genproto/blob/main/googleapis/cloud/security/privateca/v1/service.pb.go
// google/cloud/security/privateca/v1/service.proto
// go get google.golang.org/genproto/... - no extra deps besides protobuf (1.27.1)/grpc (1.40.0)

func NewGoogleCASClientRaw(capool string, creds credentials.PerRPCCredentials) (*GoogleCASClient, error) {
	caClient := &GoogleCASClient{capool: capool}
	ctx := context.Background()
	var err error
	var ol []grpc.DialOption
	ol = append(ol, grpc.WithPerRPCCredentials(creds))
	ol = append(ol, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	cc, err := grpc.DialContext(ctx, "privateca.googleapis.com:443", ol...)
	if err != nil {
		log.Printf("unable to initialize google cas caclient: %v", err)
		return nil, err
	}

	caClient.caClientRaw = privatecapb.NewCertificateAuthorityServiceClient(cc)

	if err != nil {
		log.Printf("unable to initialize google cas caclient: %v", err)
		return nil, err
	}
	return caClient, nil
}

// Extract Region from ClusterLocation
func Region(cl string) string {
	p := strings.Split(cl, "-")
	if len(p) < 3 {
		return cl
	}
	return strings.Join(p[0:2], "-")
}

func (r *GoogleCASClient) createCertReq(name string, csrPEM []byte, lifetime time.Duration) *privatecapb.CreateCertificateRequest {
	var isCA bool = false

	// We use Certificate_Config option to ensure that we only request a certificate with CAS supported extensions/usages.
	// CAS uses the PEM encoded CSR only for its public key and infers the certificate SAN (identity) of the workload through SPIFFE identity reflection
	creq := &privatecapb.CreateCertificateRequest{
		Parent: r.capool,
		// Required, [a-zA-Z0-9_-]{1,63}
		CertificateId: name,
		Certificate: &privatecapb.Certificate{
			Lifetime: durationpb.New(lifetime),
			CertificateConfig: &privatecapb.Certificate_Config{
				Config: &privatecapb.CertificateConfig{
					SubjectConfig: &privatecapb.CertificateConfig_SubjectConfig{
						Subject: &privatecapb.Subject{},
					},
					X509Config: &privatecapb.X509Parameters{
						KeyUsage: &privatecapb.KeyUsage{
							BaseKeyUsage: &privatecapb.KeyUsage_KeyUsageOptions{
								DigitalSignature: true,
								KeyEncipherment:  true,
							},
							ExtendedKeyUsage: &privatecapb.KeyUsage_ExtendedKeyUsageOptions{
								ServerAuth: true,
								ClientAuth: true,
							},
						},
						CaOptions: &privatecapb.X509Parameters_CaOptions{
							IsCa: &isCA,
						},
					},
					PublicKey: &privatecapb.PublicKey{
						Format: privatecapb.PublicKey_PEM,
						Key:    csrPEM,
					},
				},
			},
			SubjectMode: privatecapb.SubjectRequestMode_REFLECTED_SPIFFE,
		},
	}
	return creq
}

// CSR Sign calls Google CAS to sign a CSR.
func (r *GoogleCASClient) CSRSign(csrPEM []byte, certValidTTLInSec int64) ([]string, error) {
	certChain := []string{}

	//rand.Seed(time.Now().UnixNano())
	// TODO: use location, pod identity
	// `[a-zA-Z0-9_-]{1,63}`
	name := fmt.Sprintf("csr-%v", time.Now().Second())
	creq := r.createCertReq(name, csrPEM, time.Duration(certValidTTLInSec)*time.Second)

	ctx := context.Background()

	var err error
	var cresp *privatecapb.Certificate
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(creq.GetParent())))
	ctx = metadata.NewOutgoingContext(ctx, md)

	cresp, err = r.caClientRaw.CreateCertificate(ctx, creq)
	if err != nil {
		return []string{}, err
	}
	certChain = append(certChain, cresp.GetPemCertificate())
	certChain = append(certChain, cresp.GetPemCertificateChain()...)
	return certChain, nil
}

// GetRootCertBundle:  Get CA certs of the pool from Google CAS API endpoint
func (r *GoogleCASClient) GetRootCertBundle() ([]string, error) {
	var rootCertMap map[string]struct{} = make(map[string]struct{})
	var trustbundle []string = []string{}
	var err error

	ctx := context.Background()

	req := &privatecapb.FetchCaCertsRequest{
		CaPool: r.capool,
	}
	resp, err := r.caClientRaw.FetchCaCerts(ctx, req)
	if err != nil {
		log.Printf("error when getting root-certs from CAS pool: %v", err)
		return trustbundle, err
	}
	for _, certChain := range resp.CaCerts {
		certs := certChain.Certificates
		rootCert := certs[len(certs)-1]
		if _, ok := rootCertMap[rootCert]; !ok {
			rootCertMap[rootCert] = struct{}{}
		}
	}

	for rootCert := range rootCertMap {
		trustbundle = append(trustbundle, rootCert)
	}
	return trustbundle, nil
}

func (r *GoogleCASClient) Close() {
}

type WorkloadCertificateConfig struct {
	ApiVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Spec       WorkloadCertificateConfigSpec `json:"spec"`
}

type WorkloadCertificateConfigSpec struct {
	CertificateAuthorityConfig CertificateAuthorityConfig `json:"certificateAuthorityConfig"`
	ValidityDurationSeconds    int64                      `json:"validityDurationSeconds,omitempty"`
	RotationWindowPercentage   int64                      `json:"rotationWindowPercentage,omitempty"`
	KeyAlgorithm               *KeyAlgorithm              `json:"keyAlgorithm,omitempty"`
}

type CertificateAuthorityConfig struct {
	MeshCAConfig                      *MeshCAConfig                      `json:"meshCAConfig,omitempty"`
	CertificateAuthorityServiceConfig *CertificateAuthorityServiceConfig `json:"certificateAuthorityServiceConfig,omitempty"`
}

type MeshCAConfig struct {
}

type CertificateAuthorityServiceConfig struct {
	// Format: //privateca.googleapis.com/projects/PROJECT_ID/locations/SUBORDINATE_CA_LOCATION/caPools/SUBORDINATE_CA_POOL_NAME
	EndpointURI string `json:"endpointURI"`
}

type KeyAlgorithm struct {
	RSA   *RSA   `json:"rsa,omitempty"`
	ECDSA *ECDSA `json:"ecdsa,omitempty"`
}

type RSA struct {
	ModulusSize int `json:"modulusSize"`
}

type ECDSA struct {
	Curve string `json:"curve"`
}

// TrustConfig is the GKE config - when used outside GKE this is passed in the mesh-env
type TrustConfigSpec struct {
	TrustStores []TrustStore `json:"trustStores"`
}

type TrustStore struct {
	TrustDomain  string        `json:"trustDomain"`
	TrustAnchors []TrustAnchor `json:"trustAnchors,omitempty"`
}

type TrustAnchor struct {
	SPIFFETrustBundleEndpoint string `json:"spiffeTrustBundleEndpoint,omitempty"`

	// Format: //privateca.googleapis.com/projects/PROJECT_ID/locations/ROOT_CA_POOL_LOCATION/caPools/ROOT_CA_POOL_NAME
	CertificateAuthorityServiceURI string `json:"certificateAuthorityServiceURI,omitempty"`

	PEMCertificate string `json:"pemCertificate,omitempty"`
}
