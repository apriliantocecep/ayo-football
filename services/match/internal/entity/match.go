package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Match struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Date       time.Time      `json:"date"`
	Venue      string         `json:"venue"`
	HomeTeamID uuid.UUID      `gorm:"type:uuid;index" json:"home_team_id"`
	AwayTeamID uuid.UUID      `gorm:"type:uuid;index" json:"away_team_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
