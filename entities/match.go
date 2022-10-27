package entities

import "time"

// Match is the representation entity for matches
type Match struct {
	ID            uint       `json:"-"`
	Slug          *string    `json:"slug"`
	Datetime      *time.Time `json:"datetime"`
	Score         *string    `json:"score"`
	HomeTeamScore *string    `json:"home_team_score"`
	AwayTeamScore *string    `json:"away_team_score"`
	HomeTeam      Team       `json:"home_team"`
	AwayTeam      Team       `json:"away_team"`
	Videos        []Video    `json:"videos"`
}
