package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

// TODO: use otelhttptrace to get httptrace (low level client traces)

func initOTel(ctx context.Context, res string) (func(), error) {
	r := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceNameKey.String(res))

	var exp trace.SpanExporter
	var err error
	//_, shutdown, err := texporter.InstallNewPipeline(
	//	[]texporter.Option {
	//		// optional exporter options
	//	},
	//	// This example code uses sdktrace.AlwaysSample sampler to sample all traces.
	//	// In a production environment or high QPS setup please use ProbabilitySampler
	//	// set at the desired probability.
	//	// Example:
	//	// sdktrace.WithConfig(sdktrace.Config {
	//	//     DefaultSampler: sdktrace.ProbabilitySampler(0.0001),
	//	// })
	//	trace.WithConfig(trace.Config{
	//		DefaultSampler: trace.AlwaysSample(),
	//	}),
	//	// other optional provider options
	//)
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
