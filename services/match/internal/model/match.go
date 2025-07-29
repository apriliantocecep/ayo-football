package model

import "time"

type CreateMatchInput struct {
	Date       time.Time
	Venue      string
	HomeTeamID string
	AwayTeamID string
}

type UpdateMatchInput struct {
	Date       time.Time
	Venue      string
	HomeTeamID string
	AwayTeamID string
}
