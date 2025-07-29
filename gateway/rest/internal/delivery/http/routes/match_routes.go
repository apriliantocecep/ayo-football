package routes

import (
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

type MatchRoutes struct {
	App            *fiber.App
	AuthMiddleware *middlewares.AuthMiddleware
	MatchHandler   *http.MatchHandler
}

func (r *MatchRoutes) Setup() {
	routes := r.App.Use(r.AuthMiddleware.BearerTokenAuthorization)

	teams := routes.Group("matches")
	teams.Post("/", r.MatchHandler.Create)
	teams.Get("/", r.MatchHandler.List)
	teams.Get("/:id", r.MatchHandler.Get)
	teams.Put("/:id", r.MatchHandler.Update)
	teams.Delete("/:id", r.MatchHandler.Delete)

}

func NewMatchRoutes(app *fiber.App, matchHandler *http.MatchHandler, authMiddleware *middlewares.AuthMiddleware) *MatchRoutes {
	return &MatchRoutes{
		App:            app,
		MatchHandler:   matchHandler,
		AuthMiddleware: authMiddleware,
	}
}
