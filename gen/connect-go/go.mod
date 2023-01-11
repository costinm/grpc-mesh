module github.com/costinm/grpc-mesh/gen/connect-go

go 1.19

replace github.com/costinm/grpc-mesh/gen/proto/go => ../proto/go

require (
	github.com/bufbuild/connect-go v1.4.1
	github.com/costinm/grpc-mesh/gen/proto/go v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf // indirect
	google.golang.org/grpc v1.45.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
