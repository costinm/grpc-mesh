package xds

import (
	"context"
	"testing"

	"github.com/costinm/grpc-mesh/pkg/istioca"
	"github.com/costinm/meshauth"
	"github.com/costinm/mk8s/k8s"
)

var ctx = context.Background()

func TestCert(t *testing.T) {
	k, err := k8s.New(ctx, nil)

	ma, err := meshauth.FromEnv(context.Background(), nil, "")
	if err != nil {
		t.Fatal(err)
	}

	ma.AuthProviders["istio-ca"] = &meshauth.AudienceOverrideTokenSource{TokenSource: k.Default, Audience: "istio-ca"}

	// CA is only exposed on 15012 - we need a HA-PROXY version with ztunnel.
	dest := &meshauth.Dest{
		Addr: "127.0.0.1:15012",
		TokenSource: "istio-ca",
		// L4Secure: true,
	}

	_, certs, err := istioca.GetCertIstio(ctx, ma, dest, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(certs)
}

func TestUXDC(t *testing.T) {
	k, err := k8s.New(ctx, nil)

	ma, err := meshauth.FromEnv(context.Background(), nil, "")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: replace with MDS
	ma.AuthProviders["istio-ca"] = &meshauth.AudienceOverrideTokenSource{TokenSource: k.Default, Audience: "istio-ca"}

	// CA is only exposed on 15012 - we need a HA-PROXY version with ztunnel.
	//dest := &meshauth.Dest{
	//	Addr: "127.0.0.1:15010",
	//	TokenSource: "istio-ca",
	//	L4Secure: true,
	//}
	x := NewXDS(&XDSConfig{
		Namespace:  "istio-system",
		Workload:   "dev",
		XDSHeaders: nil,
		IP:         "1.2.3.4",
		XDS:        "127.0.0.1:15010",
		Context:    context.Background(),
	})

	//go func() {
	//	err := x.RunDelta("ptr")
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}()

	go func() {
		err = x.RunFull("cluster")
		if err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		err = x.RunDelta("cluster")
		//err = x.RunDelta("ptr")
		if err != nil {
			t.Fatal(err)
		}

	}()

	select{}

}

