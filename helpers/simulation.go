package helpers

import (
	"fmt"
	"insider-league/models"
	"math"
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

// CalculatePredictions calculates championship chances for each team based on their points and current week
func CalculatePredictions(teams []models.Team, currentWeek int) []models.Prediction {
	numTeams := len(teams)
	if numTeams == 0 {
		return []models.Prediction{}
	}

	// Total weeks in a season: each team plays every other team twice (home and away)
	totalWeeks := 0
	if numTeams > 1 {
		totalWeeks = 2 * (numTeams - 1)
	}

	remainingWeeks := totalWeeks - currentWeek

	// If season is complete or past completion, winner is determined
	if remainingWeeks <= 0 {
		predictions := make([]models.Prediction, numTeams)
		for i, team := range teams {
			if i == 0 { // Leader
				predictions[i] = models.Prediction{TeamName: team.Name, Chance: "100.0%"}
			} else { // Others
				predictions[i] = models.Prediction{TeamName: team.Name, Chance: "0.0%"}
			}
		}
		return predictions
	}

	// Calculate total points and check if all points are equal
	totalPoints := 0
	allPointsEqual := true
	firstTeamPoints := -1
	if numTeams > 0 {
		firstTeamPoints = teams[0].Stats.Points
	}

	for _, team := range teams {
		totalPoints += team.Stats.Points
		if team.Stats.Points != firstTeamPoints {
			allPointsEqual = false
		}
	}

	// If all teams have identical points (or no points scored yet), distribute chances equally.
	if totalPoints == 0 || allPointsEqual {
		predictions := make([]models.Prediction, numTeams)
		equalChance := 100.0 / float64(numTeams)
		for i, team := range teams {
			predictions[i] = models.Prediction{
				TeamName: team.Name,
				Chance:   fmt.Sprintf("%.1f%%", equalChance),
			}
		}
		return predictions
	}

	// Calculate progress through the season (0 to 1)
	progress := 0.0
	if totalWeeks > 0 {
		progress = math.Min(float64(currentWeek)/float64(totalWeeks), 0.99999)
	}

	// Calculate exponent that increases as season progresses
	// This will make predictions more skewed towards teams with more points
	exponent := 1.0 + progress

	// Calculate raw probabilities with exponential adjustment
	rawProbabilities := make([]float64, numTeams)
	totalRawProbability := 0.0

	for i, team := range teams {
		if totalPoints > 0 {
			// Add 1 to the points to handle the case of zero points
			rawProbabilities[i] = math.Pow(float64(team.Stats.Points+1), exponent)
			totalRawProbability += rawProbabilities[i]
		}
	}

	// Normalize probabilities to sum to 100%
	predictions := make([]models.Prediction, numTeams)
	for i, team := range teams {
		percentage := 0.0
		if totalRawProbability > 0 {
			percentage = (rawProbabilities[i] / totalRawProbability) * 100.0
		}
		predictions[i] = models.Prediction{
			TeamName: team.Name,
			Chance:   fmt.Sprintf("%.1f%%", percentage),
		}
	}

	return predictions
}
