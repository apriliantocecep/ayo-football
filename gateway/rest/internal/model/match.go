package model

type CreateMatchRequest struct {
	Date       string `json:"date" validate:"required,datetime=2006-01-02 15:04:05"`
	Venue      string `json:"venue" validate:"required"`
	HomeTeamID string `json:"home_team_id" validate:"required"`
	AwayTeamID string `json:"away_team_id" validate:"required"`
}

type UpdateMatchRequest = CreateMatchRequest

type CreateGoalRequest struct {
	//MatchID  string `json:"match_id" validate:"required"`
	PlayerID string `json:"player_id" validate:"required"`
	ScoredAt string `json:"scored_at" validate:"required,datetime=2006-01-02 15:04:05"`
}

type UpdateGoalRequest = CreateGoalRequest
