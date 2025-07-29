package routes

import (
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

type PlayerRoutes struct {
	App            *fiber.App
	AuthMiddleware *middlewares.AuthMiddleware
	PlayerHandler  *http.PlayerHandler
}

func (r *PlayerRoutes) Setup() {
	routes := r.App.Use(r.AuthMiddleware.BearerTokenAuthorization)

	players := routes.Group("players")
	players.Post("/", r.PlayerHandler.CreatePlayer)
	players.Get("/", r.PlayerHandler.ListPlayerByTeam)
	players.Get("/:id", r.PlayerHandler.GetPlayer)
	players.Put("/:id", r.PlayerHandler.UpdatePlayer)
	players.Delete("/:id", r.PlayerHandler.DeletePlayer)

}

func NewPlayerRoutes(app *fiber.App, playerHandler *http.PlayerHandler, authMiddleware *middlewares.AuthMiddleware) *PlayerRoutes {
	return &PlayerRoutes{
		App:            app,
		PlayerHandler:  playerHandler,
		AuthMiddleware: authMiddleware,
	}
}
