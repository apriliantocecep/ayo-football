package routes

import (
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type AuthRoutes struct {
	App         *fiber.App
	AuthHandler *http.AuthHandler
}

func (r *AuthRoutes) Setup() {
	auth := r.App.Group("auth")
	auth.Post("/login", r.AuthHandler.Login)
	auth.Post("/register", r.AuthHandler.Register)
}

func NewAuthRoutes(app *fiber.App, authHandler *http.AuthHandler) *AuthRoutes {
	return &AuthRoutes{App: app, AuthHandler: authHandler}
}
