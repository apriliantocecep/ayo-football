package utils

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ExtractTokenFromHeader(c *fiber.Ctx) (string, error) {
	authorization := c.Get("Authorization")

	if authorization == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return token[1], nil
}
