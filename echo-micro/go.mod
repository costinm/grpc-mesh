module github.com/costinm/grpc-mesh/echo-micro

go 1.17

replace github.com/costinm/grpc-mesh/bootstrap => ../bootstrap

replace github.com/costinm/grpc-mesh/telemetry/otelgrpc => ../telemetry/otelgrpc

replace github.com/costinm/grpc-mesh/telemetry/logs/zap => ../telemetry/logs/zap

replace github.com/costinm/grpc-mesh/gen/proto/go => ../gen/proto/go

replace google.golang.org/grpc => ../../grpc

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.0
	github.com/costinm/grpc-mesh/bootstrap v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/gen/proto/go v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/telemetry/logs/zap v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/telemetry/otelgrpc v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/prometheus/client_golang v1.12.1
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0
	go.opentelemetry.io/contrib/instrumentation/host v0.31.0
	go.opentelemetry.io/contrib/zpages v0.31.0
	go.opentelemetry.io/otel v1.6.1
	go.opentelemetry.io/otel/exporters/prometheus v0.28.0
	go.opentelemetry.io/otel/metric v0.28.0
	go.opentelemetry.io/otel/sdk v1.6.1
	go.opentelemetry.io/otel/sdk/metric v0.28.0
	go.uber.org/zap v1.21.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.45.0
)

require (
	cloud.google.com/go/compute v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20220112060520-0fa49ea1db0c // indirect
	github.com/envoyproxy/go-control-plane v0.10.2-0.20220119144013-5ce31a86d483 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.2 // indirect
	github.com/go-kit/log v0.1.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/shirou/gopsutil/v3 v3.22.2 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/otel/trace v1.6.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
