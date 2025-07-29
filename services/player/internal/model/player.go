package model

type CreatePlayerInput struct {
	TeamID     string
	Name       string
	Height     float32
	Weight     float32
	Position   string
	BackNumber int32
}

type UpdatePlayerInput struct {
	TeamID     string
	Name       string
	Height     float32
	Weight     float32
	Position   string
	BackNumber int32
}
