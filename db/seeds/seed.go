package seeds

import (
	"insider-league/models"
	"log"

	"gorm.io/gorm"
)

// Load seeds the database with initial teams and matches
func Load(db *gorm.DB) error {
	// Check if teams already exist
	var count int64
	if err := db.Model(&models.Team{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already seeded.")
		return nil
	}

	// Create teams
	teams := []models.Team{
		{Name: "Chelsea", Strength: 85, Stats: models.Stats{}},
		{Name: "Arsenal", Strength: 87, Stats: models.Stats{}},
		{Name: "Manchester City", Strength: 94, Stats: models.Stats{}},
		{Name: "Liverpool", Strength: 92, Stats: models.Stats{}},
	}

	// Save teams to database
	if err := db.Create(&teams).Error; err != nil {
		return err
	}

	// Generate matches
	matches := generateFixtures(teams)

	// Save matches to database
	if err := db.Create(&matches).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully with 4 teams and 12 matches.")
	return nil
}

// generateFixtures creates all matches for the season
func generateFixtures(teams []models.Team) []models.Match {
	var matches []models.Match
	numTeams := len(teams)

	// For odd number of teams, we'll need to add a "bye" team
	// This ensures we always have an even number of teams to schedule
	if numTeams%2 != 0 {
		log.Printf("Warning: Odd number of teams (%d). Adding a bye team.", numTeams)
		// Add a dummy team for bye weeks
		teams = append(teams, models.Team{Name: "BYE", Strength: 0})
		numTeams = len(teams)
	}

	weeksPerHalf := numTeams - 1
	totalWeeks := weeksPerHalf * 2

	// Create a copy of teams slice to rotate
	teamRotation := make([]models.Team, len(teams))
	copy(teamRotation, teams)

	// First half of the season
	for week := 1; week <= weeksPerHalf; week++ {
		// For each week, create matches between teams at opposite ends of the rotation
		// For n teams, we create n/2 matches per week
		for i := 0; i < numTeams/2; i++ {
			homeTeam := teamRotation[i]
			awayTeam := teamRotation[numTeams-1-i]

			// Skip matches involving the BYE team
			if homeTeam.Name == "BYE" || awayTeam.Name == "BYE" {
				continue
			}

			// Home match
			matches = append(matches, models.Match{
				Week:          week,
				HomeTeamID:    homeTeam.ID,
				AwayTeamID:    awayTeam.ID,
				HomeTeamScore: 0,
				AwayTeamScore: 0,
				IsPlayed:      false,
			})

			// Away match (in second half)
			matches = append(matches, models.Match{
				Week:          week + weeksPerHalf,
				HomeTeamID:    awayTeam.ID,
				AwayTeamID:    homeTeam.ID,
				HomeTeamScore: 0,
				AwayTeamScore: 0,
				IsPlayed:      false,
			})
		}

		// Rotate teams for next week (keep first team fixed, rotate others)
		// Move the second team to the end, shift others up
		secondTeam := teamRotation[1]
		copy(teamRotation[1:], teamRotation[2:])
		teamRotation[len(teamRotation)-1] = secondTeam
	}

	log.Printf("Generated %d matches over %d weeks", len(matches), totalWeeks)
	return matches
}
