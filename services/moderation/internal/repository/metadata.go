package repository

import (
	"errors"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MetadataRepository struct {
}

func (m *MetadataRepository) FindByArticleId(db *gorm.DB, articleId string) (*entity.Metadata, error) {
	var metadata entity.Metadata
	if err := db.Where(&entity.Metadata{ArticleId: articleId}).First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}
	return &metadata, nil
}

func (m *MetadataRepository) Update(db *gorm.DB, metadata *entity.Metadata) error {
	return db.Save(metadata).Error
}

func (m *MetadataRepository) Create(db *gorm.DB, metadata *entity.Metadata) (uuid.UUID, error) {
	result := db.Create(&metadata)
	return metadata.ID, result.Error
}

func NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{}
}
