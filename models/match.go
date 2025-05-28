package models

// Match represents a football match in the league
type Match struct {
	ID            uint `json:"id" db:"id" gorm:"primaryKey"`
	Week          int  `json:"week" db:"week"`
	HomeTeamID    uint `json:"homeTeamId" db:"home_team_id"`
	AwayTeamID    uint `json:"awayTeamId" db:"away_team_id"`
	HomeTeamScore int  `json:"homeTeamScore" db:"home_team_score"`
	AwayTeamScore int  `json:"awayTeamScore" db:"away_team_score"`
	IsPlayed      bool `json:"isPlayed" db:"is_played"`
}
