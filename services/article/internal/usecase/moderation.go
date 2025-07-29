package usecase

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/article/internal/model"
	"github.com/apriliantocecep/ayo-football/services/article/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModerationUseCase struct {
	DB                *mongo.Client
	ArticleRepository *repository.ArticleRepository
}

func (u *ModerationUseCase) PublishArticle(ctx context.Context, request *model.PublishArticleRequest) (*model.PublishArticleResponse, error) {
	articleStatus := "draft"
	if request.ModerationStatus == "accepted" {
		articleStatus = "published"
	}
	if request.ModerationStatus == "rejected" {
		articleStatus = "rejected"
	}
	updateData := bson.M{
		"status": articleStatus,
	}
	err := u.ArticleRepository.Update(request.ArticleId, updateData)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not update article status: %v", err)
	}

	response := model.PublishArticleResponse{Status: articleStatus}
	return &response, nil
}

func NewModerationUseCase(DB *mongo.Client, articleRepository *repository.ArticleRepository) *ModerationUseCase {
	return &ModerationUseCase{DB: DB, ArticleRepository: articleRepository}
}
