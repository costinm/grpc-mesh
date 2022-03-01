// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoogleCloudPlatform/cloud-run-mesh/pkg/mesh"
	"github.com/costinm/grpc-mesh/echo-micro/proto"
	"github.com/costinm/grpc-mesh/echo-micro/server"
	"go.opencensus.io/zpages"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	"log"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/runmetrics"

	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/reflection"
)

// TODO: use otelhttptrace to get httptrace (low level client traces)

func HttpRoundTripper() func(transport http.RoundTripper) http.RoundTripper {
	return func(transport http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(transport)
	}
}

func initOTel(ctx context.Context, kr *mesh.KRun) (func(), error) {

	var exp trace.SpanExporter
	var err error
	exp, err = stdouttrace.New(
		stdouttrace.WithWriter(os.Stderr),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		//stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}

	r := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceNameKey.String(kr.Name))
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// =========== Metrics
	var exporter metric.Exporter
	exporter, err = stdoutmetric.New(
		//stdoutmetric.WithPrettyPrint(),
		stdoutmetric.WithWriter(os.Stderr))

	if err != nil {
		log.Fatalf("creating stdoutmetric exporter: %v", err)
	}

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithInexpensiveDistribution(),
			exporter,
			processor.WithMemory(true),
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(3*time.Second),
		controller.WithResource(r),
		// WithResource, WithCollectPeriod, WithPushTimeout
	)

	if err = pusher.Start(ctx); err != nil {
		log.Fatalf("starting push controller: %v", err)
	}

	global.SetMeterProvider(pusher)

	// Global instrumentations
	if err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	); err != nil {
		log.Fatalln("failed to start runtime instrumentation:", err)
	}
	// Host telemetry -
	host.Start()

	// End telemetry magic
	return func() {
		if err := pusher.Stop(ctx); err != nil {
			log.Fatalf("stopping push controller: %v", err)
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}, nil
	/*
		kr.TransportWrapper = func(transport http.RoundTripper) http.RoundTripper {
			return otelhttp.NewTransport(transport)
		}
		// Host telemetry -
		host.Start()
	*/
}

func OTELGRPCClient() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor())}
}

func Run(port string) error {
	h := &server.EchoGrpcHandler{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	creds := insecure.NewCredentials()

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}

	err = runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
		Prefix:       "echo/",
	})
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	proto.RegisterEchoTestServiceServer(grpcServer, h)
	admin.Register(grpcServer)
	reflection.Register(grpcServer)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	// Status
	mux := &http.ServeMux{}
	zpages.Handle(mux, "/debug")

	http.ListenAndServe("127.0.0.1:8081", mux)
	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	return nil
}

// Most minimal gRPC based server, for estimating binary size overhead.
//
// - 0.8M for a min go program
// - 4.7M for an echo using HTTP.
// - 9M - this server, only plain gRPC
// - 20M - same app, but proxyless gRPC
// - 22M - plus opencensus, prom, zpages, reflection
func main() {
	err := Run(":8080")
	if err != nil {
		fmt.Println("Error ", err)
	}
}
