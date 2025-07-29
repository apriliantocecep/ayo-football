package middlewares

import (
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/services/auth/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	AuthServiceClient *grpc_client.AuthServiceClient
}

func (m *AuthMiddleware) BearerTokenAuthorization(c *fiber.Ctx) error {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		return err
	}

	req := pb.ValidateTokenRequest{Token: token}
	res, err := m.AuthServiceClient.Client.ValidateToken(c.UserContext(), &req)
	if err != nil {
		return utils.HandleGrpcError(err)
		//return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Locals("user_id", res.UserId)

	return c.Next()
}

func NewAuthMiddleware(authServiceClient *grpc_client.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{AuthServiceClient: authServiceClient}
}
