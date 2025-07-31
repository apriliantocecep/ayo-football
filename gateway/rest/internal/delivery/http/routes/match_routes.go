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

	matches := routes.Group("matches")
	matches.Post("/", r.MatchHandler.Create)
	matches.Get("/", r.MatchHandler.List)
	matches.Get("/:id", r.MatchHandler.Get)
	matches.Put("/:id", r.MatchHandler.Update)
	matches.Delete("/:id", r.MatchHandler.Delete)
	// goals
	matches.Post("/:id/goals", r.MatchHandler.CreateGoal)
	matches.Get("/:id/goals/:goalId", r.MatchHandler.GetGoal)
	matches.Put("/:id/goals/:goalId", r.MatchHandler.UpdateGoal)
	matches.Delete("/:id/goals/:goalId", r.MatchHandler.DeleteGoal)

}

func NewMatchRoutes(app *fiber.App, matchHandler *http.MatchHandler, authMiddleware *middlewares.AuthMiddleware) *MatchRoutes {
	return &MatchRoutes{
		App:            app,
		MatchHandler:   matchHandler,
		AuthMiddleware: authMiddleware,
	}
}
