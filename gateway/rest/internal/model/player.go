package model

import (
	"github.com/apriliantocecep/ayo-football/services/player/pkg/pb"
	"github.com/go-playground/validator/v10"
	"strings"
)

type CreatePlayerRequest struct {
	TeamID     string  `json:"team_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Height     float32 `json:"height" validate:"required"`
	Weight     float32 `json:"weight" validate:"required"`
	Position   string  `json:"position" validate:"required,position_enum"`
	BackNumber int32   `json:"back_number" validate:"required,min=1,max=100"`
}

type UpdatePlayerRequest = CreatePlayerRequest

var ValidPositions = map[string]pb.Position{
	"PENYERANG":      pb.Position_PENYERANG,
	"GELANDANG":      pb.Position_GELANDANG,
	"BERTAHAN":       pb.Position_BERTAHAN,
	"PENJAGA_GAWANG": pb.Position_PENJAGA_GAWANG,
}

func PositionEnumValidation(fl validator.FieldLevel) bool {
	_, ok := ValidPositions[strings.ToUpper(fl.Field().String())]
	return ok
}

type PlayerResource struct {
	ID         string  `json:"id,omitempty"`
	TeamID     string  `json:"team_id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Height     float32 `json:"height,omitempty"`
	Weight     float32 `json:"weight,omitempty"`
	Position   string  `json:"position,omitempty"`
	BackNumber int32   `json:"back_number,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

func PlayerToResponse(player *pb.Player) *PlayerResource {
	return &PlayerResource{
		ID:         player.Id,
		TeamID:     player.TeamId,
		Name:       player.Name,
		Height:     player.Height,
		Weight:     player.Weight,
		Position:   player.Position.String(),
		BackNumber: player.BackNumber,
		CreatedAt:  player.CreatedAt,
		UpdatedAt:  player.UpdatedAt,
	}
}
