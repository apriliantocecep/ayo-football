package usecase

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/model"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MetadataUseCase struct {
	DB                 *gorm.DB
	MetadataRepository *repository.MetadataRepository
}

func (u *MetadataUseCase) Save(ctx context.Context, request *model.MetadataRequest) (*model.MetadataResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	metadata := entity.Metadata{
		ArticleId:        request.ArticleId,
		Title:            request.Title,
		Author:           request.Author,
		ModerationStatus: "pending",
	}
	articleUuid, err := u.MetadataRepository.Create(u.DB, &metadata)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create metadata article")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create metadata article")
	}

	response := model.MetadataResponse{MetadataId: articleUuid.String()}
	return &response, nil
}

func NewMetadataUseCase(DB *gorm.DB, metadataRepository *repository.MetadataRepository) *MetadataUseCase {
	return &MetadataUseCase{DB: DB, MetadataRepository: metadataRepository}
}
