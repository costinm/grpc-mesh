module github.com/costinm/grpc-mesh/echo-micro

go 1.16

replace github.com/costinm/grpc-mesh/bootstrap => ../bootstrap

replace github.com/costinm/grpc-mesh/telemetry/otelgrpc => ../telemetry/otelgrpc

replace github.com/costinm/grpc-mesh/telemetry/logs/zap => ../telemetry/logs/zap

replace github.com/costinm/grpc-mesh/gen/proto/go => ../gen/proto/go

replace github.com/costinm/grpc-mesh => ../

replace google.golang.org/grpc => ../../grpc

replace github.com/GoogleCloudPlatform/cloud-run-mesh => ../../cloud-run-mesh

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.0
	github.com/GoogleCloudPlatform/cloud-run-mesh v0.0.0-20220128230121-cac57262761b
	github.com/costinm/grpc-mesh/bootstrap v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/gen/proto/go v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/telemetry/logs/zap v0.0.0-00010101000000-000000000000
	github.com/costinm/grpc-mesh/telemetry/otelgrpc v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0 // indirect
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
	cloud.google.com/go/container v1.2.0 // indirect
	cloud.google.com/go/security v1.3.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
)
