package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http/middlewares"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http/routes"
	"github.com/apriliantocecep/posfin-blog/shared"
	sharedlib "github.com/apriliantocecep/posfin-blog/shared/lib"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// context
	ctx := context.Background()

	// Fiber instance
	app := fiber.New(fiber.Config{
		AppName:      "Posfin Service Gateway",
		ErrorHandler: newErrorHandler(),
		IdleTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Prefork:      true,
	})
	app.Use(recover.New())
	//app.Use(logger.New(logger.Config{
	//	TimeZone: "Asia/Jakarta",
	//}))

	// start vault client
	vaultClient := shared.NewVaultClient()

	// otel
	otelSDK := sharedlib.NewOtelSDK(ctx, vaultClient, "gateway-srv")
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
	tracer := tp.Tracer("gateway.main.tracer")

	prop := otelSDK.NewPropagator()
	metricExporter, err := otelSDK.OTLPMetricExporter()
	if err != nil {
		log.Fatalf("failed to create OTLP metric exporter: %v", err)
	}
	meterProvider := otelSDK.NewMeterProvider(otelmetric.NewPeriodicReader(metricExporter))
	defer func() {
		if err = meterProvider.Shutdown(ctx); err != nil {
			log.Printf("error shutting down meter provider: %v", err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	// gRPC Clients
	authServiceClient := grpc_client.NewAuthServiceClient(vaultClient, tp, prop, meterProvider)
	defer func(Conn *grpc.ClientConn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("closing connection to auth service error: %v", err)
		}
	}(authServiceClient.Conn)

	articleServiceClient := grpc_client.NewArticleServiceClient(vaultClient)
	defer func(Conn *grpc.ClientConn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("closing connection to article service error: %v", err)
		}
	}(articleServiceClient.Conn)

	// validator instance
	newValidator := validator.New()

	// http handlers
	authHandler := http.NewAuthHandler(authServiceClient, newValidator, tracer)
	articleHandler := http.NewArticleHandler(newValidator, articleServiceClient, authServiceClient)

	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authServiceClient)

	// routes
	authRoutes := routes.NewAuthRoutes(app, authHandler)
	articleRoutes := routes.NewArticleRoutes(app, articleHandler, authMiddleware)
	authRoutes.Setup()
	articleRoutes.Setup()

	// listener
	secret := utils.GetVaultSecretConfig(vaultClient)
	portStr := secret["GATEWAY_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalln("GATEWAY_PORT is not set")
	}
	port := utils.ParsePort(portStr.(string))

	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("failed to start gateway server: %+v\n", err)
	}
}

func newErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errs := make(map[string]string)
			for _, ve := range validationErrors {
				errs[ve.Field()] = utils.ValidationErrorMessage(ve)
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": errs,
			})
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
