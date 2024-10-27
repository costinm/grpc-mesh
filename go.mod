module github.com/costinm/grpc-mesh

go 1.22.2

replace github.com/costinm/grpc-mesh/gen/proto => ./gen/proto

replace github.com/costinm/meshauth => ../meshauth

require (
	github.com/bufbuild/connect-go v1.10.0
	github.com/costinm/grpc-mesh/gen/connect-go v0.0.0-20240605142039-df0642efa7d0
	github.com/costinm/grpc-mesh/gen/proto v0.0.0-20240605230615-7973e12eb107
	github.com/costinm/meshauth v0.0.0-20240803190121-2a6dfc0e888a
	github.com/costinm/mk8s/k8s v0.0.0-20240627225616-dad11d70021b
	github.com/costinm/ssh-mesh v0.0.0-20240729162626-2efe1420fa7f
	github.com/golang/protobuf v1.5.4
	go.opentelemetry.io/proto/otlp v1.3.1
	golang.org/x/crypto v0.26.0
	golang.org/x/net v0.28.0
	google.golang.org/genproto v0.0.0-20240725223205-93522f1f2a9f
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	sigs.k8s.io/yaml v1.4.0
)

require (
	cloud.google.com/go/longrunning v0.5.11 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/ginkgo/v2 v2.19.1 // indirect
	github.com/onsi/gomega v1.34.0 // indirect
	github.com/pkg/sftp v1.13.6 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.30.3 // indirect
	k8s.io/apimachinery v0.30.3 // indirect
	k8s.io/client-go v0.30.3 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240726031636-6f6746feab9c // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)
