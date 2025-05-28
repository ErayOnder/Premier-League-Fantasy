package models

// Stats represents the statistical data for a team
type Stats struct {
	Points         int `json:"points" gorm:"column:points"`
	GoalsFor       int `json:"goals_for" gorm:"column:goals_for"`
	GoalsAgainst   int `json:"goals_against" gorm:"column:goals_against"`
	GoalDifference int `json:"goal_difference" gorm:"column:goal_difference"`
	Wins           int `json:"wins" gorm:"column:wins"`
	Draws          int `json:"draws" gorm:"column:draws"`
	Losses         int `json:"losses" gorm:"column:losses"`
}

// Team represents a football team in the league
type Team struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Strength int    `json:"strength"`
	Stats    Stats  `json:"stats" gorm:"embedded"`
}
