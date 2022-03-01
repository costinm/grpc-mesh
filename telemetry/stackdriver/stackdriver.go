package stackdriver

import (
	"context"
	"time"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.7.0"

	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// Broken:
	// github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric imports
	//	go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount: module go.opentelemetry.io/otel/sdk/metric@latest found (v0.26.0), but does not contain package go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount
	sdmetrics "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"

	sdtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

// Helper to init stackdriver using Otel.
func InitStackdriver(ctx context.Context, projectId string,
		creds credentials.PerRPCCredentials, res string) error {
	// k8s based GSA federated access and ID token provider
	//tokenProvider, _ := sts.NewSTS(kr)
	//tokenProvider.MDPSA = true
	//tokenProvider.UseAccessToken = true

	var exp trace.SpanExporter
	var err error
	exp, err = sdtrace.New(sdtrace.WithProjectID(projectId),
		sdtrace.WithTraceClientOptions([]option.ClientOption{
			option.WithGRPCDialOption(grpc.WithPerRPCCredentials(creds)),
			option.WithQuotaProject(projectId),
		}))
	if err != nil {
		return err
	}

	r := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceNameKey.String(res))

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
		propagation.Baggage{}))

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

	// Metrics
	//var exporter metric.Exporter

	//pusher
	_, err = sdmetrics.InstallNewPipeline([]sdmetrics.Option{sdmetrics.WithProjectID(projectId),
		sdmetrics.WithMonitoringClientOptions(
			option.WithGRPCDialOption(grpc.WithPerRPCCredentials(creds)),
			option.WithQuotaProject(projectId))},
		//controller.WithExporter(exporter),
		controller.WithCollectPeriod(3*time.Second),
		controller.WithResource(r),
	)

	//exporter, err = mexporter.NewRawExporter(mexporter.WithProjectID(projectId),
	//	mexporter.WithMonitoringClientOptions(
	//		option.WithGRPCDialOption(grpc.WithPerRPCCredentials(creds)),
	//		option.WithQuotaProject(projectId)))
	//
	//pusher := controller.New(
	//	processor.NewFactory(
	//		simple.NewWithInexpensiveDistribution(),
	//		exporter,
	//		processor.WithMemory(true),
	//	),
	//	controller.WithExporter(exporter),
	//	controller.WithCollectPeriod(3*time.Second),
	//	controller.WithResource(r),
	//	// WithResource, WithCollectPeriod, WithPushTimeout
	//)
	//
	//if err = pusher.Start(ctx); err != nil {
	//	return err
	//}
	//
	//global.SetMeterProvider(pusher)

	return nil
}
