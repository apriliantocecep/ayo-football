package http

import (
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/posfin-blog/gateway/rest/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/article/pkg/pb"
	authpb "github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	Validate             *validator.Validate
	ArticleServiceClient *grpc_client.ArticleServiceClient
	AuthServiceClient    *grpc_client.AuthServiceClient
}

func (h *ArticleHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	request := new(model.ArticleRequest)
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

	userRes, err := h.AuthServiceClient.Client.GetUserById(c.UserContext(), &authpb.GetUserByIdRequest{Id: userID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	req := pb.SubmitArticleRequest{
		Title:       request.Title,
		Author:      userRes.GetName(),
		HtmlContent: request.Content,
		UserId:      userID,
	}

	res, err := h.ArticleServiceClient.Client.SubmitArticle(c.UserContext(), &req)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{
		"article_id": res.GetArticleId(),
		"status":     res.GetStatus(),
	})
}

func NewArticleHandler(validate *validator.Validate, articleServiceClient *grpc_client.ArticleServiceClient, authServiceClient *grpc_client.AuthServiceClient) *ArticleHandler {
	return &ArticleHandler{
		Validate:             validate,
		ArticleServiceClient: articleServiceClient,
		AuthServiceClient:    authServiceClient,
	}
}
