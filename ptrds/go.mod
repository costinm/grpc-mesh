module github.com/costinm/grpc-mesh/ptrds

go 1.19

replace github.com/costinm/grpc-mesh/gen/proto => ../gen/proto

replace github.com/costinm/grpc-mesh/gen/connect-go => ../gen/connect-go

require (
	github.com/bufbuild/connect-go v1.4.1
	github.com/costinm/grpc-mesh/gen/connect-go v0.0.0-00010101000000-000000000000
)

require (
	github.com/costinm/grpc-mesh/gen/proto v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
