module github.com/costinm/grpc-mesh/gen/connect-go

go 1.19

replace github.com/costinm/grpc-mesh/gen/proto => ../proto

require (
	github.com/bufbuild/connect-go v1.4.1
	github.com/costinm/grpc-mesh/gen/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
