package models

// Prediction represents a team's chance of winning a match
type Prediction struct {
	TeamName string `json:"teamName"`
	Chance   int    `json:"chance"`
}
