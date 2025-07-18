package usecase

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/entity"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/repository"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type ArticleUseCase struct {
	DB                      *mongo.Client
	ArticleRepository       *repository.ArticleRepository
	ArticleCreatedPublisher *messaging.ArticlePublisher
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

	// publish to broker
	event := sharedmodel.ArticleEvent{
		ID:     articleID,
		Title:  request.Title,
		Author: request.Author,
	}
	err = u.ArticleCreatedPublisher.Publish(&event)
	if err != nil {
		log.Printf("failed publish article created event : %+v", err)
		return nil, status.Errorf(codes.Aborted, "failed to publish article metadata")
	}

	response := model.ArticleResponse{ArticleId: articleID, Status: articleStatus}

	return &response, nil
}

func NewArticleUseCase(DB *mongo.Client, articleRepository *repository.ArticleRepository, articleCreatedPublisher *messaging.ArticlePublisher) *ArticleUseCase {
	return &ArticleUseCase{
		DB:                      DB,
		ArticleRepository:       articleRepository,
		ArticleCreatedPublisher: articleCreatedPublisher,
	}
}
