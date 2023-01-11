# Telemetry helpers and plugins 

Helpers to configure otel tracing and monitoring for gRPC

## Existing solutions

### Istio abstraction

istio.io/pkg/monitoring
 - wrapper on top of opencensus
 - defines an OnRecordHook - useful for testing but adds overhead
 - ocensus is deprecated and has its own overhead

### OpenCensus instrumentation 

- (deprecated) - prom exporter (159), stackdriver (292), ocgrpc (434)
- used in GoogleCloudPlatform examples

###  OpenTelemetry 

(not ready yet)

### Prometheus

github.com/grpc-ecosystem/go-grpc-middleware/providers/openmetrics/v2 (1) (using directly github.com/prometheus/client_golang/prometheus) 
  and the old github.com/grpc-ecosystem/go-grpc-prometheus (1600 imports)

- tiller
- https://grafana.com/grafana/dashboards/14765 (139 downloads)
- last commit in 2021/03 - and 2021/08 for the openmetrics in go-grpc-middleware

Metrics:

- grpc_server_started_total{"grpc_type", "grpc_service", "grpc_method"}
- grpc_server_handled_total{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}
- grpc_server_msg_received_total{"grpc_type", "grpc_service", "grpc_method"}
- grpc_server_msg_sent_total{"grpc_type", "grpc_service", "grpc_method"}
- same for client
