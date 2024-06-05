module github.com/costinm/grpc-mesh/gen/connect

go 1.22

replace github.com/costinm/grpc-mesh/gen/proto => ../proto

require (
	github.com/bufbuild/connect-go v1.4.1
	github.com/costinm/grpc-mesh/gen/proto v0.0.0-00010101000000-000000000000
	github.com/costinm/ugate/gen/proto v0.0.0-20221024013023-789def6d5dde
	go.opentelemetry.io/proto/otlp v1.2.0
	google.golang.org/grpc v1.63.0
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
