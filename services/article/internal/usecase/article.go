package usecase

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/entity"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArticleUseCase struct {
	DB                *mongo.Client
	ArticleRepository *repository.ArticleRepository
}

func (u *ArticleUseCase) Insert(ctx context.Context, request *model.ArticleRequest) (*model.ArticleResponse, error) {
	articleStatus := "draft"
	article := entity.Article{
		Content: request.Content,
		UserId:  request.UserId,
		Status:  articleStatus,
	}

	articleID, err := u.ArticleRepository.Create(&article)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create article content")
	}

	response := model.ArticleResponse{ArticleId: articleID, Status: articleStatus}

	return &response, nil
}

func NewArticleUseCase(DB *mongo.Client, articleRepository *repository.ArticleRepository) *ArticleUseCase {
	return &ArticleUseCase{DB: DB, ArticleRepository: articleRepository}
}
