package routes

import (
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

type TeamRoutes struct {
	App            *fiber.App
	AuthMiddleware *middlewares.AuthMiddleware
	TeamHandler    *http.TeamHandler
}

func (r *TeamRoutes) Setup() {
	routes := r.App.Use(r.AuthMiddleware.BearerTokenAuthorization)

	teams := routes.Group("teams")
	teams.Post("/", r.TeamHandler.CreateTeam)
	teams.Get("/", r.TeamHandler.ListTeams)
	teams.Get("/:id", r.TeamHandler.GetTeam)
	teams.Put("/:id", r.TeamHandler.UpdateTeam)
	teams.Delete("/:id", r.TeamHandler.DeleteTeam)

}

func NewTeamRoutes(app *fiber.App, teamHandler *http.TeamHandler, authMiddleware *middlewares.AuthMiddleware) *TeamRoutes {
	return &TeamRoutes{
		App:            app,
		TeamHandler:    teamHandler,
		AuthMiddleware: authMiddleware,
	}
}
