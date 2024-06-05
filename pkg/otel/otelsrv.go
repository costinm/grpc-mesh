package otel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	connect_go "github.com/bufbuild/connect-go"
	logsconnect "github.com/costinm/grpc-mesh/gen/connect-go/opentelemetry/proto/collector/logs/v1/v1connect"
	metricsconnect "github.com/costinm/grpc-mesh/gen/connect-go/opentelemetry/proto/collector/metrics/v1/v1connect"
	traceconnect "github.com/costinm/grpc-mesh/gen/connect-go/opentelemetry/proto/collector/trace/v1/v1connect"
	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	metricsv1 "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

type Config struct {

}

// Receives otel spans, metrics and logs using gRPC protocol.
// Currently, it just logs to slog - the intent is to provide a circular
// buffer and allow getting the data for debug and tools.
//
// Similar to `otel-cli server json --endpoint http://0.0.0.0:4318  --stdout` - but for all
//
// Test with
//   otel-cli span --name "my-span" --start 1622522888 --end 1622523000 --service "my-service" --kind "client" --verbose --endpoint localhost:4317
//   otel-cli span --name "my-span" --start 1622522888 --end 1622523000 --service "my-service" --kind "client" --verbose --endpoint http://localhost:4318
type Otel struct {
	*Config
	Traces *OTelSvc
	Logs *OTelSvcLogs
	Metrics *OTelSvcMetrics
}

func NewOtel(cfg *Config) *Otel {
	o := &Otel{
		Config: cfg,
		Logs: &OTelSvcLogs{},
		Metrics: &OTelSvcMetrics{},
		Traces: &OTelSvc{},
	}

	return o
}

func (ot *Otel) Register(mux *http.ServeMux) {
	mux.Handle(metricsconnect.NewMetricsServiceHandler(ot.Metrics))
	mux.Handle(traceconnect.NewTraceServiceHandler(ot.Traces))
	mux.Handle(logsconnect.NewLogsServiceHandler(ot.Logs))

	mux.HandleFunc("/v1/traces", ot.HandleHTTP)
	mux.HandleFunc("/v1/logs", ot.HandleHTTP)
	mux.HandleFunc("/v1/metrics", ot.HandleHTTP)
}

func (o *Otel) HandleHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("Error while reading request body: %s", err)
		}

		var msg proto.Message
		var res proto.Message
		if strings.HasSuffix(req.URL.Path, "/v1/traces") {
			msg = &tracev1.ExportTraceServiceRequest{}
			res = &tracev1.ExportTraceServiceResponse{}
		} else if strings.HasSuffix(req.URL.Path, "/v1/logs") {
			msg = &logsv1.ExportLogsServiceRequest{}
			res = &logsv1.ExportLogsServiceResponse{}
		} else if strings.HasSuffix(req.URL.Path, "/v1/metrics") {
			msg = &metricsv1.ExportMetricsServiceRequest{}
			res = &metricsv1.ExportMetricsServiceResponse{}
		} else {
			writer.WriteHeader(http.StatusNotAcceptable)
			return
		}

		var resBytes []byte
		ct := req.Header.Get("Content-Type")
		writer.Header().Set("Content-Type", ct)
		switch ct {
		case "application/x-protobuf":
			proto.Unmarshal(data, msg)
			resBytes, _  = proto.Marshal(res)
		case "application/json":
			json.Unmarshal(data, msg)
			resBytes, _  = json.Marshal(res)
		default:
			writer.WriteHeader(http.StatusNotAcceptable)
			return
		}

		writer.WriteHeader(200)

		// Notes:
		// span kind client=1 server=2
		// spanid and traceid encoded as bytes instead of hex

		msgb, _ := yaml.Marshal(msg)
		fmt.Println("---\n# \n" + string(msgb))

		writer.Write(resBytes)
	} else {
		writer.WriteHeader(404)
		writer.Write([]byte("Not found"))
	}

}

type OTelSvc struct {
	traceconnect.UnimplementedTraceServiceHandler
}

func (*OTelSvc) Export(ctx context.Context, req *connect_go.Request[tracev1.ExportTraceServiceRequest]) (*connect_go.Response[tracev1.ExportTraceServiceResponse], error) {
	msgb, _ := yaml.Marshal(req.Msg)
	fmt.Println("---\n# \n" + string(msgb))
	return connect_go.NewResponse(&tracev1.ExportTraceServiceResponse{}), nil
}

type OTelSvcLogs struct {
	logsconnect.UnimplementedLogsServiceHandler
}

func (*OTelSvcLogs) Export(ctx context.Context, req *connect_go.Request[logsv1.ExportLogsServiceRequest]) (*connect_go.Response[logsv1.ExportLogsServiceResponse], error) {
	msgb, _ := yaml.Marshal(req.Msg)
	fmt.Println("---\n# \n" + string(msgb))
	return connect_go.NewResponse(&logsv1.ExportLogsServiceResponse{}), nil
}

type OTelSvcMetrics struct {
	logsconnect.UnimplementedLogsServiceHandler
}

func (*OTelSvcMetrics) Export(ctx context.Context, req *connect_go.Request[metricsv1.ExportMetricsServiceRequest]) (*connect_go.Response[metricsv1.ExportMetricsServiceResponse], error) {
	msgb, _ := yaml.Marshal(req.Msg)
	fmt.Println("---\n# \n" + string(msgb))
	return connect_go.NewResponse(&metricsv1.ExportMetricsServiceResponse{}), nil
}
