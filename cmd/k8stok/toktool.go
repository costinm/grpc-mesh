package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/gcp"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/mesh"
	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/sts"
)

var (
	aud   = flag.String("aud", "istio-ca", "Audience, if empty, the k8s default token is returned")
	gcpSA = flag.String("gsa", "", "Google service account to impersonate")

	fed = flag.Bool("fed", false, "Return the federated token")

	namespace = flag.String("n", "default", "Namespace")
	//sa        = flag.String("sa", "default", "Service account")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	mesh.Debug = false

	kr := mesh.New()
	kr.SkipSaveCerts = true
	kr.Namespace = *namespace
	err := gcp.InitGCP(ctx, kr)

	err = kr.LoadConfig(ctx)
	if err != nil {
		panic(err)
	}

	tokenProvider, err := sts.NewSTS(kr)

	if *gcpSA == "" && !*fed {
		tokenProvider.K8S = true
		tokenProvider.AudOverride = *aud
	} else if *fed {
		tokenProvider.UseAccessToken = true
	} else {
		if *gcpSA == "default" {
			tokenProvider.MDPSA = true
		} else {
			tokenProvider.GSA = *gcpSA
		}
	}

	t, err := tokenProvider.GetRequestMetadata(ctx, *aud)
	if err != nil {
		log.Fatal("Failed to get token", err)
	}
	at := t["authorization"]
	att := at[7:]
	fmt.Println(att)
}
