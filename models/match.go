package models

// Match represents a football match in the league
type Match struct {
	ID            int `json:"id" db:"id"`
	Week          int `json:"week" db:"week"`
	HomeTeamID    int `json:"homeTeamId" db:"home_team_id"`
	AwayTeamID    int `json:"awayTeamId" db:"away_team_id"`
	HomeTeamScore int `json:"homeTeamScore" db:"home_team_score"`
	AwayTeamScore int `json:"awayTeamScore" db:"away_team_score"`
}
