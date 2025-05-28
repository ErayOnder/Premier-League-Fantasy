package models

// Prediction represents a team's chance of winning the championship
type Prediction struct {
	TeamName string `json:"teamName"`
	Chance   string `json:"chance"`
}
