package helpers

import (
	"fmt"
	"insider-league/models"
	"math/rand"
)

// SimulateMatchScore generates a random match score based on the relative strengths of the teams.
func SimulateMatchScore(homeTeamStrength, awayTeamStrength int) (homeGoals, awayGoals int) {
	// Calculate total strength for probability distribution
	totalStrength := homeTeamStrength + awayTeamStrength

	// Maximum number of goals that can be scored by a team in a match
	maxGoals := 5

	// Simulate home team goals
	for i := range maxGoals {
		// Probability of scoring decreases with each goal
		probability := float64(homeTeamStrength) / float64(totalStrength) * (1.0 - float64(i)/float64(maxGoals))
		if rand.Float64() < probability {
			homeGoals++
		} else {
			break
		}
	}

	// Simulate away team goals
	for i := range maxGoals {
		probability := float64(awayTeamStrength) / float64(totalStrength) * (1.0 - float64(i)/float64(maxGoals)) * 0.9
		if rand.Float64() < probability {
			awayGoals++
		} else {
			break
		}
	}

	return homeGoals, awayGoals
}

// CalculatePredictions calculates championship chances for each team based on their points
func CalculatePredictions(teams []models.Team) []models.Prediction {
	// Calculate total points
	totalPoints := 0
	for _, team := range teams {
		totalPoints += team.Stats.Points
	}

	// If no points have been scored yet, return empty predictions
	if totalPoints == 0 {
		return []models.Prediction{}
	}

	// Calculate predictions for each team
	predictions := make([]models.Prediction, len(teams))
	for i, team := range teams {
		percentage := float64(team.Stats.Points) * 100 / float64(totalPoints)
		predictions[i] = models.Prediction{
			TeamName: team.Name,
			Chance:   fmt.Sprintf("%.1f%%", percentage),
		}
	}

	return predictions
}
