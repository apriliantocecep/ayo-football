package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/delivery/grpc_server"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/repository"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared"
	sharedlib "github.com/apriliantocecep/posfin-blog/shared/lib"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	capi "github.com/hashicorp/consul/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/stats/opentelemetry"
	"log"
	"net"
	"os"
)

func main() {
	// setup vars
	ctx := context.Background()
	serviceName := "auth-service-cluster"
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatalln("PORT is not set")
	}
	port := utils.ParsePort(portStr)

	url := os.Getenv("SERVICE_URL")
	if url == "" {
		log.Fatalln("SERVICE_URL is not set")
	}
	address := url

	// vault client
	vaultClient := shared.NewVaultClient()
	//secret := utils.GetVaultSecretConfig(vaultClient)

	// otel
	otelSDK := sharedlib.NewOtelSDK(ctx, vaultClient, "auth-srv")
	spanExporter, err := otelSDK.OTLPSpanExporter()
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}
	tp := otelSDK.NewTraceProvider(spanExporter)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	//tracer := tp.Tracer("gateway.main.tracer")
	prop := otelSDK.NewPropagator()

	// consul client
	consul := sharedlib.NewConsulClient(vaultClient)

	// rabbitmq client
	rabbitMQClient := sharedlib.NewRabbitMQClient(vaultClient)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing rabbitmq: %v", err)
		}
	}(rabbitMQClient.Conn)

	// setup publisher
	userCreatedPublisher := messaging.NewUserPublisher(rabbitMQClient.Conn, "user_created", "user_created")

	// dependencies
	database := config.NewDatabase(vaultClient)
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}(database.Conn)
	jwt := config.NewJwt(vaultClient)
	userRepository := repository.NewUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepository, jwt, database.DB, userCreatedPublisher)

	// grpc server
	so := opentelemetry.ServerOption(opentelemetry.Options{
		//MetricsOptions: opentelemetry.MetricsOptions{
		//	MeterProvider: meterProvider,
		//	// These are example experimental gRPC metrics, which are disabled
		//	// by default and must be explicitly enabled. For the full,
		//	// up-to-date list of metrics, see:
		//	// https://grpc.io/docs/guides/opentelemetry-metrics/#instruments
		//	Metrics: opentelemetry.DefaultMetrics().Add(
		//		"grpc.server.call.started",
		//		"grpc.server.call.duration",
		//	),
		//},
		TraceOptions: oteltracing.TraceOptions{TracerProvider: tp, TextMapPropagator: prop}},
	)

	srv := grpc_server.NewAuthServer(userUseCase)
	s := grpc.NewServer(so)
	pb.RegisterAuthServiceServer(s, srv)

	// listener
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", listen.Addr())

	// register service to consul
	serviceRegisteredID := fmt.Sprintf("auth-service-%d", port)
	tags := []string{
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s.local`)", serviceName, serviceName),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=grpc", serviceName),
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.scheme=h2c", serviceName),
	}
	consul.ServiceRegister(serviceRegisteredID, serviceName, address, port, &capi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Interval:                       "10s",
		Timeout:                        "2s",
		DeregisterCriticalServiceAfter: "1m",
	}, tags)

	// health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
