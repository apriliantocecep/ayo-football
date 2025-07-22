package grpc_client

import (
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
	"log"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
	Conn   *grpc.ClientConn
}

func NewAuthServiceClient(vaultClient *shared.VaultClient, traceProvider trace.TracerProvider, textMapPropagator propagation.TextMapPropagator, meterProvider *otelmetric.MeterProvider) *AuthServiceClient {
	secret := utils.GetVaultSecretConfig(vaultClient)

	proxyUrl := secret["AUTH_SERVICE_PROXY"]
	if proxyUrl == nil || proxyUrl == "" {
		log.Fatalln("AUTH_SERVICE_PROXY is not set")
	}
	grpcProxyUrl := secret["TRAEFIK_GRPC_PROXY_URL"]
	if grpcProxyUrl == nil || grpcProxyUrl == "" {
		log.Fatalln("TRAEFIK_GRPC_PROXY_URL is not set")
	}

	// Dial Options
	do := opentelemetry.DialOption(opentelemetry.Options{
		//MetricsOptions: opentelemetry.MetricsOptions{
		//	MeterProvider: meterProvider,
		//	// These are example experimental gRPC metrics, which are disabled
		//	// by default and must be explicitly enabled. For the full,
		//	// up-to-date list of metrics, see:
		//	// https://grpc.io/docs/guides/opentelemetry-metrics/#instruments
		//	Metrics: opentelemetry.DefaultMetrics().Add(
		//		"grpc.client.attempt.started",
		//		"grpc.client.attempt.duration",
		//	),
		//},
		TraceOptions: oteltracing.TraceOptions{TracerProvider: traceProvider, TextMapPropagator: textMapPropagator},
	})

	conn, err := grpc.NewClient(
		grpcProxyUrl.(string),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithAuthority(proxyUrl.(string)),
		do,
	)
	if err != nil {
		log.Fatalf("did not connect to auth service: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		Client: client,
		Conn:   conn,
	}
}
