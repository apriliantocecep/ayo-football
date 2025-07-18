package repository

import (
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetadataRepository struct {
}

func (m *MetadataRepository) Create(db *gorm.DB, metadata *entity.Metadata) (uuid.UUID, error) {
	result := db.Create(&metadata)
	return metadata.ID, result.Error
}

func NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{}
}
