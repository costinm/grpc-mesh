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

package xdsc

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/gcp"
	_ "github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/gcp"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/k8s"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/mesh"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/sts"
)

// TestSTS uses a k8s connection and env to locate the mesh, and tests the token generation.
// The tokens are needed to connect to TD, MCP, ASM, CA.
func TestSTS(t *testing.T) {
	kr := mesh.New()

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	err := gcp.InitGCP(ctx, kr)
	if err != nil {
		log.Fatal("Failed to find K8S ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}

	err = kr.LoadConfig(ctx)
	if err != nil {
		t.Skip("Failed to connect to GKE, missing kubeconfig ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}

	if kr.ProjectNumber == "" {
		t.Skip("Skipping STS test, PROJECT_NUMBER required")
	}
	masterT, err := kr.GetToken(ctx, kr.TrustDomain)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(sts.TokenPayload(masterT))

	s, err := sts.NewSTS(&sts.AuthConfig{
		ProjectNumber:  kr.ProjectNumber,
		TrustDomain:    kr.TrustDomain,
		ClusterAddress: kr.ClusterAddress,
		TokenSource:    kr,
	})
	if err != nil {
		t.Fatal(err)
	}

	f, err := s.TokenFederated(ctx, masterT)
	if err != nil {
		t.Fatal(err)
	}

	a, err := s.TokenAccess(ctx, f, "")
	if err != nil {
		t.Fatal(err)
	}

	a, err = s.TokenAccess(ctx, f, "https://foo.bar")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(sts.TokenPayload(a))
}

func TestP4SA(t *testing.T) {
	// Get an access token for MCP
	ctx, cf := context.WithTimeout(context.Background(), 100*time.Second)
	defer cf()

	kr := mesh.New()

	kc := &k8s.K8S{Mesh: kr}
	err := kc.K8SClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	kr.Cfg = kc
	kr.TokenProvider = kc

	err = kr.LoadConfig(ctx)
	if err != nil {
		t.Skip("Failed to connect to GKE, missing kubeconfig ", time.Since(kr.StartTime), kr, os.Environ(), err)
	}

	sts, err := sts.NewSTS(&sts.AuthConfig{
		ProjectNumber:  kr.ProjectNumber,
		TrustDomain:    kr.TrustDomain,
		ClusterAddress: kr.ClusterAddress,
		TokenSource:    kr,
	})

	masterT, err := kr.GetToken(ctx, kr.TrustDomain)
	if err != nil {
		t.Fatal(err)
	}

	f, err := sts.TokenFederated(ctx, masterT)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sts.TokenAccess(ctx, f, "")
	if err != nil {
		t.Fatal(err)
	}
	//log.Println(a)
}
