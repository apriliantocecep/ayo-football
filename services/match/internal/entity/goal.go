package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Goal struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	MatchID   uuid.UUID      `gorm:"type:uuid" json:"match_id"`
	PlayerID  uuid.UUID      `gorm:"type:uuid" json:"player_id"`
	ScoredAt  time.Time      `json:"date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
