package main

import (
	"errors"
	"fmt"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http/middlewares"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http/routes"
	"github.com/apriliantocecep/posfin-blog/shared"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// Fiber instance
	app := fiber.New(fiber.Config{
		AppName:      "Posfin Service Gateway",
		ErrorHandler: newErrorHandler(),
	})
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Jakarta",
	}))

	// start vault client
	vaultClient := shared.NewVaultClient()

	// gRPC Clients
	authServiceClient := grpc_client.NewAuthServiceClient(vaultClient)
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
	authHandler := http.NewAuthHandler(authServiceClient, newValidator)
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
	portStr := secret["GATEWAY_PORT"].(string)
	port := utils.ParsePort(portStr)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start gateway server: %v", err)
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
