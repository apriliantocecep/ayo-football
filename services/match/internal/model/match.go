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

type CreateGoalInput struct {
	MatchID  string
	PlayerID string
	ScoredAt time.Time
}

type UpdateGoalInput struct {
	MatchID  string
	PlayerID string
	ScoredAt time.Time
}
