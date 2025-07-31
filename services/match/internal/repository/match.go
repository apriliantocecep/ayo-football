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
	CreateGoal(ctx context.Context, db *gorm.DB, goal *entity.Goal) error
	GetGoalByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Goal, error)
	UpdateGoal(ctx context.Context, db *gorm.DB, goal *entity.Goal) error
	DeleteGoal(ctx context.Context, db *gorm.DB, id uuid.UUID) error
}

type MatchRepository struct{}

func (r *MatchRepository) CreateGoal(ctx context.Context, db *gorm.DB, goal *entity.Goal) error {
	goal.ID = uuid.New()
	return db.WithContext(ctx).Create(goal).Error
}

func (r *MatchRepository) GetGoalByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Goal, error) {
	var goal entity.Goal
	if err := db.WithContext(ctx).First(&goal, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *MatchRepository) UpdateGoal(ctx context.Context, db *gorm.DB, goal *entity.Goal) error {
	return db.WithContext(ctx).Save(goal).Error
}

func (r *MatchRepository) DeleteGoal(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return db.WithContext(ctx).Delete(&entity.Goal{}, "id = ?", id).Error
}

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
