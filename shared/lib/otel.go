package lib

import (
	"context"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Otel struct {
	Context     context.Context
	VaultClient *shared.VaultClient
	ServiceName string
}

func (o *Otel) NewMeterProvider(exporter otelmetric.Reader) *otelmetric.MeterProvider {
	// Ensure default SDK resources and the required service name are set.
	res, err := o.NewResource()
	if err != nil {
		log.Fatalf("failed to create otel meter resource: %v", err)
	}
	meterProvider := otelmetric.NewMeterProvider(
		otelmetric.WithResource(res),
		otelmetric.WithReader(exporter),
	)
	return meterProvider
}

func (o *Otel) NewTraceProvider(exporter oteltrace.SpanExporter) *oteltrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	res, err := o.NewResource()
	if err != nil {
		log.Fatalf("failed to create otel trace resource: %v", err)
	}

	tp := oteltrace.NewTracerProvider(
		oteltrace.WithBatcher(exporter),
		oteltrace.WithResource(res),
	)

	return tp
}

func (o *Otel) NewResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(o.ServiceName)))
}

func (o *Otel) NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (o *Otel) OTLPSpanExporter() (oteltrace.SpanExporter, error) {
	secret := utils.GetVaultSecretConfig(o.VaultClient)
	endpoint := secret["OTEL_COLLECTOR_ENDPOINT_HTTP"]
	if endpoint == nil || endpoint == "" {
		log.Fatalln("OTEL_COLLECTOR_ENDPOINT_HTTP is not set")
	}

	return otlptracehttp.New(o.Context,
		otlptracehttp.WithEndpoint(endpoint.(string)),
		otlptracehttp.WithInsecure(),
	)
}

func (o *Otel) OTLPGRPCExporter() (oteltrace.SpanExporter, error) {
	secret := utils.GetVaultSecretConfig(o.VaultClient)
	endpoint := secret["OTEL_COLLECTOR_ENDPOINT_GRPC"]
	if endpoint == nil || endpoint == "" {
		log.Fatalln("OTEL_COLLECTOR_ENDPOINT_GRPC is not set")
	}

	conn, err := grpc.NewClient(
		endpoint.(string),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err)
	}

	return otlptracegrpc.New(o.Context, otlptracegrpc.WithGRPCConn(conn))
}

func (o *Otel) OTLPMetricExporter() (otelmetric.Exporter, error) {
	secret := utils.GetVaultSecretConfig(o.VaultClient)
	endpoint := secret["OTEL_COLLECTOR_ENDPOINT_HTTP"]
	if endpoint == nil || endpoint == "" {
		log.Fatalln("OTEL_COLLECTOR_ENDPOINT_HTTP is not set")
	}

	return otlpmetrichttp.New(o.Context,
		otlpmetrichttp.WithEndpoint(endpoint.(string)),
		otlpmetrichttp.WithInsecure(),
	)
}

func NewOtelSDK(context context.Context, vaultClient *shared.VaultClient, serviceName string) *Otel {
	return &Otel{Context: context, VaultClient: vaultClient, ServiceName: serviceName}
}
