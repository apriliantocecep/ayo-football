package model

type CreateTeamInput struct {
	Name      string
	Logo      string
	FoundedAt int
	Address   string
	City      string
}

type UpdateTeamInput struct {
	Name      string
	Logo      string
	FoundedAt int
	Address   string
	City      string
}
