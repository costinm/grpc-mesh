module github.com/costinm/grpc-mesh/cas

go 1.18

replace github.com/costinm/grpc-mesh/gen/grpc-go => ../gen/grpc-go

require (
	cloud.google.com/go/security v1.11.0
	github.com/costinm/grpc-mesh/gen/grpc-go v0.0.0-20230111230202-a4e8d3c436dd
	github.com/costinm/meshauth v0.0.0-20230114063639-beded2028d83
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	google.golang.org/grpc v1.52.0
	google.golang.org/protobuf v1.28.1
	sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go/longrunning v0.4.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230113154510-dbe35b8444a5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
