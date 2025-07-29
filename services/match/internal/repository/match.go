package repository

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/match/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMatchRepository interface {
	Create(ctx context.Context, db *gorm.DB, match *entity.Match) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Match, error)
	Update(ctx context.Context, db *gorm.DB, match *entity.Match) error
	Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error
	List(ctx context.Context, db *gorm.DB, offset, limit int) ([]*entity.Match, error)
}

type MatchRepository struct{}

func (r *MatchRepository) Create(ctx context.Context, db *gorm.DB, match *entity.Match) error {
	match.ID = uuid.New()
	return db.WithContext(ctx).Create(match).Error
}

func (r *MatchRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Match, error) {
	var match entity.Match
	if err := db.WithContext(ctx).First(&match, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) Update(ctx context.Context, db *gorm.DB, match *entity.Match) error {
	return db.WithContext(ctx).Save(match).Error
}

func (r *MatchRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return db.WithContext(ctx).Delete(&entity.Match{}, "id = ?", id).Error
}

func (r *MatchRepository) List(ctx context.Context, db *gorm.DB, offset, limit int) ([]*entity.Match, error) {
	var matches []*entity.Match
	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

var _ IMatchRepository = &MatchRepository{}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{}
}
