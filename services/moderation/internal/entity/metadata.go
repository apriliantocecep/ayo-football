package entity

import (
	"github.com/google/uuid"
	"time"
)

type Metadata struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	ArticleId        string    `json:"article_id" gorm:"not null"`
	Title            string    `json:"title" gorm:"not null"`
	Author           string    `json:"author"`
	ModerationStatus string    `json:"moderation_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
