package repository

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/player/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerRepository struct{}

func (r *PlayerRepository) Create(ctx context.Context, db *gorm.DB, player *entity.Player) error {
	player.ID = uuid.New()
	return db.WithContext(ctx).Create(player).Error
}

func (r *PlayerRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Player, error) {
	var player entity.Player
	if err := db.WithContext(ctx).First(&player, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) Update(ctx context.Context, db *gorm.DB, player *entity.Player) error {
	return db.WithContext(ctx).Save(player).Error
}

func (r *PlayerRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return db.WithContext(ctx).Delete(&entity.Player{}, "id = ?", id).Error
}

func (r *PlayerRepository) ListByTeamID(ctx context.Context, db *gorm.DB, teamID uuid.UUID, offset, limit int) ([]*entity.Player, error) {
	var players []*entity.Player
	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Where("team_id = ?", teamID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func (r *PlayerRepository) IsBackNumberUsed(ctx context.Context, db *gorm.DB, teamID uuid.UUID, backNumber int32, excludeID string) (bool, error) {
	var count int64
	query := db.WithContext(ctx).Model(&entity.Player{}).
		Where("team_id = ? AND back_number = ?", teamID, backNumber)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewPlayerRepository() *PlayerRepository {
	return &PlayerRepository{}
}
