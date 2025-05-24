package models

// Stats represents the statistical data for a team
type Stats struct {
	Points       int `json:"points" db:"points"`
	GoalsFor     int `json:"goals_for" db:"goals_for"`
	GoalsAgainst int `json:"goals_against" db:"goals_against"`
	Wins         int `json:"wins" db:"wins"`
	Draws        int `json:"draws" db:"draws"`
	Losses       int `json:"losses" db:"losses"`
}

// Team represents a football team in the league
type Team struct {
	ID       uint   `json:"id" db:"id" gorm:"primaryKey"`
	Name     string `json:"name" db:"name"`
	Strength int    `json:"strength" db:"strength"`
	Stats    Stats  `json:"stats" db:"stats"`
}
