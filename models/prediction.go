package models

// Prediction represents a team's chance of winning the championship
type Prediction struct {
	TeamName string `json:"teamName"`
	Chance   int    `json:"chance"`
}
