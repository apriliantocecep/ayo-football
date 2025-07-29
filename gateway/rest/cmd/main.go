package main

import (
	"errors"
	"fmt"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http/middlewares"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http/routes"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// Fiber instance
	app := fiber.New(fiber.Config{
		AppName:      "Ayo football Service Gateway",
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

	// gRPC Clients
	authServiceClient := grpc_client.NewAuthServiceClient(vaultClient)
	defer func(Conn *grpc.ClientConn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("closing connection to auth service error: %v", err)
		}
	}(authServiceClient.Conn)
	teamServiceClient := grpc_client.NewTeamServiceClient(vaultClient)
	defer func(Conn *grpc.ClientConn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("closing connection to team service error: %v", err)
		}
	}(teamServiceClient.Conn)

	// validator instance
	newValidator := validator.New()

	// http handlers
	authHandler := http.NewAuthHandler(authServiceClient, newValidator)
	teamHandler := http.NewTeamHandler(newValidator, teamServiceClient)

	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authServiceClient)

	// routes
	authRoutes := routes.NewAuthRoutes(app, authHandler)
	authRoutes.Setup()
	teamRoutes := routes.NewTeamRoutes(app, teamHandler, authMiddleware)
	teamRoutes.Setup()

	// listener
	secret := utils.GetVaultSecretConfig(vaultClient)
	portStr := secret["GATEWAY_PORT"]
	if portStr == nil || portStr == "" {
		log.Fatalln("GATEWAY_PORT is not set")
	}
	port := utils.ParsePort(portStr.(string))

	err := app.Listen(fmt.Sprintf(":%d", port))
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
