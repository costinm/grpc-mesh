# istio-grpc
Helpers and samples for using gRPC XDS with Istio. Includes a dependency-free XDS client.

## Logging

Istio uses Zap, with a wrapper that adds a lot of dependencies.

GRPC also has a wrapper (grpc/grpclog/LoggerV2 interface), defaulting to the built-in go log, configured via env:
GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99

It is possible to use Zap as grpc logger and dynamically configure it:
```go
    // 	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
    //  "go.uber.org/zap"
    
    
    // '4' skips the adapters, shows the actual code logging the info.
	// 3  works for most core - but fails for 'prefixLogger' - used by XDS
	zl, _ := zap.NewDevelopment(zap.AddCallerSkip(4))
	// TODO: use zap.observer, we can use the collected info for assertions.
	grpc_zap.ReplaceGrpcLoggerV2WithVerbosity(zl, 99)

```

Grpc sources also include grpc/grpclog/glogger, which depends on the dep-free github.com/golang/glog
The logger can also intercept the standard log, and add file:line and send to files, and allows
configuring log level per source file/module. Unfortunately this is only configurable via flags,
so not very useful in tests or libraries.
