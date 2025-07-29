package repository

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/team/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct{}

func (r *TeamRepository) Create(ctx context.Context, db *gorm.DB, team *entity.Team) error {
	team.ID = uuid.New()
	return db.Create(team).Error
}

func (r *TeamRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	if err := db.First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) Update(ctx context.Context, db *gorm.DB, team *entity.Team) error {
	//return db.Model(&entity.Team{}).Where("id = ?", team.ID).Updates(team).Error
	return db.Save(team).Error
}

func (r *TeamRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&entity.Team{}, "id = ?", id).Error
}

func (r *TeamRepository) ListAll(ctx context.Context, db *gorm.DB, offset, limit int) ([]entity.Team, error) {
	var teams []entity.Team
	if err := db.Offset(offset).Limit(limit).Order("created_at desc").Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}
