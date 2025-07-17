package routes

import (
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

type ArticleRoutes struct {
	App            *fiber.App
	AuthMiddleware *middlewares.AuthMiddleware
	ArticleHandler *http.ArticleHandler
}

func (r *ArticleRoutes) Setup() {
	routes := r.App.Use(r.AuthMiddleware.BearerTokenAuthorization)

	auth := routes.Group("articles")
	auth.Post("/", r.ArticleHandler.Create)
}

func NewArticleRoutes(app *fiber.App, articleHandler *http.ArticleHandler, authMiddleware *middlewares.AuthMiddleware) *ArticleRoutes {
	return &ArticleRoutes{
		App:            app,
		ArticleHandler: articleHandler,
		AuthMiddleware: authMiddleware,
	}
}
