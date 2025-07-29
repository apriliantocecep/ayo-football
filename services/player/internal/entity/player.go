package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Player struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeamID     uuid.UUID      `gorm:"type:uuid;index" json:"team_id"`
	Name       string         `json:"name"`
	Height     float32        `json:"height"`   // cm
	Weight     float32        `json:"weight"`   // kg
	Position   string         `json:"position"` // enum-like
	BackNumber int32          `json:"back_number"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
