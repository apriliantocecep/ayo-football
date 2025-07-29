package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `json:"name"`
	Logo      string         `json:"logo"`
	FoundedAt int            `json:"founded_at"`
	Address   string         `json:"address"`
	City      string         `json:"city"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
