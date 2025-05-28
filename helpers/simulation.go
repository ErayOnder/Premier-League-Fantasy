package helpers

import (
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
