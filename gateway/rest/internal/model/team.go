package model

type CreateTeamRequest struct {
	Name      string `json:"name" validate:"required"`
	Logo      string `json:"logo" validate:"required,url"`
	FoundedAt int32  `json:"founded_at" validate:"required"`
	Address   string `json:"address" validate:"required"`
	City      string `json:"city" validate:"required"`
}

type UpdateTeamRequest = CreateTeamRequest
