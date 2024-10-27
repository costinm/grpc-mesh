# Grpc-mesh

Minimal implementations of mesh related services, for testing and small binaries.

Like other implementations attempting to reduce dependencies, it has a fork of the protos - which are supposed to 
stay stable and backwards compatible. 

To keep size and dependencies low - it is using the connect implementation that uses the go native http stack. It is
also a way to check how well it interoperates with Istiod and other services and to learn about it.

The proto directory has an assortment of protos I consider useful for mesh - envoy XDS, few gRPC and fortio 
test endpoints - but also OTel. Not all are currently used.

## Deps

Using connect-go for minimal binary size and to use the native http stack - it is integrated with other mesh services.

The gen/proto is kept separated - it doesn't depend on connect or any other package.

Meshauth is used to bootstrap, no deps.

Connect-go depends on bufbuild/connecto-go.
