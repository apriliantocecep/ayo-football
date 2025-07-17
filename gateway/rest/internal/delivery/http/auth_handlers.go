package http

import (
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validate          *validator.Validate
	AuthServiceClient *grpc_client.AuthServiceClient
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	err := c.BodyParser(request)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: err.Error(),
		}
	}

	err = h.Validate.Struct(request)
	if err != nil {
		return err
	}

	req := pb.RegisterRequest{
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
	}

	res, err := h.AuthServiceClient.Client.Register(c.UserContext(), &req)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{
		"user_id":  res.GetUserId(),
		"username": res.GetUsername(),
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	err := c.BodyParser(request)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: err.Error(),
		}
	}

	err = h.Validate.Struct(request)
	if err != nil {
		return err
	}

	req := pb.LoginRequest{
		Identifier: request.Identity,
		Password:   request.Password,
	}

	res, err := h.AuthServiceClient.Client.Login(c.UserContext(), &req)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{
		"access_token": res.GetAccessToken(),
		"expire_at":    res.GetExpiresAt(),
	})
}

func NewAuthHandler(authServiceClient *grpc_client.AuthServiceClient, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{AuthServiceClient: authServiceClient, Validate: validator}
}
