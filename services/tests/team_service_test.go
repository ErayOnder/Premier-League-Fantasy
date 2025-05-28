package tests

import (
	repomocks "insider-league/mocks/repository"
	"insider-league/models"
	"insider-league/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTeamService_UpdateMatchStats(t *testing.T) {
	tests := []struct {
		name              string
		homeGoals         int
		awayGoals         int
		revert            bool
		expectedHomeStats models.Stats
		expectedAwayStats models.Stats
		description       string
	}{
		{
			name:      "Home team win (3-1)",
			homeGoals: 3,
			awayGoals: 1,
			revert:    false,
			expectedHomeStats: models.Stats{
				Points:         13, // 10 + 3
				Wins:           3,  // 2 + 1
				Draws:          1,  // unchanged
				Losses:         1,  // unchanged
				GoalsFor:       8,  // 5 + 3
				GoalsAgainst:   4,  // 3 + 1
				GoalDifference: 4,  // 8 - 4
			},
			expectedAwayStats: models.Stats{
				Points:         6,  // 6 + 0
				Wins:           2,  // unchanged
				Draws:          0,  // unchanged
				Losses:         3,  // 2 + 1
				GoalsFor:       5,  // 4 + 1
				GoalsAgainst:   8,  // 5 + 3
				GoalDifference: -3, // 5 - 8
			},
			description: "Home team wins, gets 3 points and 1 win, away team gets 1 loss",
		},
		{
			name:      "Away team win (0-2)",
			homeGoals: 0,
			awayGoals: 2,
			revert:    false,
			expectedHomeStats: models.Stats{
				Points:         10, // 10 + 0
				Wins:           2,  // unchanged
				Draws:          1,  // unchanged
				Losses:         2,  // 1 + 1
				GoalsFor:       5,  // 5 + 0
				GoalsAgainst:   5,  // 3 + 2
				GoalDifference: 0,  // 5 - 5
			},
			expectedAwayStats: models.Stats{
				Points:         9, // 6 + 3
				Wins:           3, // 2 + 1
				Draws:          0, // unchanged
				Losses:         2, // unchanged
				GoalsFor:       6, // 4 + 2
				GoalsAgainst:   5, // unchanged
				GoalDifference: 1, // 6 - 5
			},
			description: "Away team wins, gets 3 points and 1 win, home team gets 1 loss",
		},
		{
			name:      "Draw (2-2)",
			homeGoals: 2,
			awayGoals: 2,
			revert:    false,
			expectedHomeStats: models.Stats{
				Points:         11, // 10 + 1
				Wins:           2,  // unchanged
				Draws:          2,  // 1 + 1
				Losses:         1,  // unchanged
				GoalsFor:       7,  // 5 + 2
				GoalsAgainst:   5,  // 3 + 2
				GoalDifference: 2,  // 7 - 5
			},
			expectedAwayStats: models.Stats{
				Points:         7,  // 6 + 1
				Wins:           2,  // unchanged
				Draws:          1,  // 0 + 1
				Losses:         2,  // unchanged
				GoalsFor:       6,  // 4 + 2
				GoalsAgainst:   7,  // 5 + 2
				GoalDifference: -1, // 6 - 7
			},
			description: "Draw, both teams get 1 point and 1 draw",
		},
		{
			name:      "Revert home team win (3-1)",
			homeGoals: 3,
			awayGoals: 1,
			revert:    true,
			expectedHomeStats: models.Stats{
				Points:         7, // 10 - 3
				Wins:           1, // 2 - 1
				Draws:          1, // unchanged
				Losses:         1, // unchanged
				GoalsFor:       2, // 5 - 3
				GoalsAgainst:   2, // 3 - 1
				GoalDifference: 0, // 2 - 2
			},
			expectedAwayStats: models.Stats{
				Points:         6, // 6 - 0
				Wins:           2, // unchanged
				Draws:          0, // unchanged
				Losses:         1, // 2 - 1
				GoalsFor:       3, // 4 - 1
				GoalsAgainst:   2, // 5 - 3
				GoalDifference: 1, // 3 - 2
			},
			description: "Reverting home team win, subtracts 3 points and 1 win from home, 1 loss from away",
		},
		{
			name:      "Revert draw (2-2)",
			homeGoals: 2,
			awayGoals: 2,
			revert:    true,
			expectedHomeStats: models.Stats{
				Points:         9, // 10 - 1
				Wins:           2, // unchanged
				Draws:          0, // 1 - 1
				Losses:         1, // unchanged
				GoalsFor:       3, // 5 - 2
				GoalsAgainst:   1, // 3 - 2
				GoalDifference: 2, // 3 - 1
			},
			expectedAwayStats: models.Stats{
				Points:         5,  // 6 - 1
				Wins:           2,  // unchanged
				Draws:          -1, // 0 - 1
				Losses:         2,  // unchanged
				GoalsFor:       2,  // 4 - 2
				GoalsAgainst:   3,  // 5 - 2
				GoalDifference: -1, // 2 - 3
			},
			description: "Reverting draw, subtracts 1 point and 1 draw from both teams",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(repomocks.MockTeamRepository)

			// Create team service with mock
			service := services.NewTeamService(mockRepo)

			// Create initial teams with some stats
			homeTeam := &models.Team{
				ID:   1,
				Name: "Home Team",
				Stats: models.Stats{
					Points:       10,
					Wins:         2,
					Draws:        1,
					Losses:       1,
					GoalsFor:     5,
					GoalsAgainst: 3,
				},
			}

			awayTeam := &models.Team{
				ID:   2,
				Name: "Away Team",
				Stats: models.Stats{
					Points:       6,
					Wins:         2,
					Draws:        0,
					Losses:       2,
					GoalsFor:     4,
					GoalsAgainst: 5,
				},
			}

			// Set up mock expectations - Update should be called for both teams
			mockRepo.On("Update", mock.MatchedBy(func(team *models.Team) bool {
				return team.ID == 1 // home team
			})).Return(nil).Once()

			mockRepo.On("Update", mock.MatchedBy(func(team *models.Team) bool {
				return team.ID == 2 // away team
			})).Return(nil).Once()

			// Call the function under test
			err := service.UpdateMatchStats(homeTeam, awayTeam, tt.homeGoals, tt.awayGoals, tt.revert)

			// Assertions
			assert.NoError(t, err, "UpdateMatchStats should not return an error")
			assert.Equal(t, tt.expectedHomeStats, homeTeam.Stats, "Home team stats should match expected values")
			assert.Equal(t, tt.expectedAwayStats, awayTeam.Stats, "Away team stats should match expected values")

			// Verify that all expected calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTeamService_GetLeagueTable(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockTeamRepository)

	// Create team service with mock
	service := services.NewTeamService(mockRepo)

	// Create unsorted teams with different stats
	unsortedTeams := []models.Team{
		{
			ID:   1,
			Name: "Team A",
			Stats: models.Stats{
				Points:       15, // Should be 2nd (same points as Team D, but lower goal difference)
				GoalsFor:     20,
				GoalsAgainst: 10, // Goal difference: +10
			},
		},
		{
			ID:   2,
			Name: "Team B",
			Stats: models.Stats{
				Points:       12, // Should be 3rd
				GoalsFor:     18,
				GoalsAgainst: 8, // Goal difference: +10
			},
		},
		{
			ID:   3,
			Name: "Team C",
			Stats: models.Stats{
				Points:       9, // Should be 4th (lowest points)
				GoalsFor:     15,
				GoalsAgainst: 12, // Goal difference: +3
			},
		},
		{
			ID:   4,
			Name: "Team D",
			Stats: models.Stats{
				Points:       15, // Should be 1st (same points as Team A, but higher goal difference)
				GoalsFor:     25,
				GoalsAgainst: 10, // Goal difference: +15
			},
		},
	}

	// Set up mock expectation
	mockRepo.On("GetAll").Return(unsortedTeams, nil).Once()

	// Call the function under test
	sortedTeams, err := service.GetLeagueTable()

	// Assertions
	assert.NoError(t, err, "GetLeagueTable should not return an error")
	assert.Len(t, sortedTeams, 4, "Should return all 4 teams")

	// Verify correct sorting order
	expectedOrder := []struct {
		name     string
		points   int
		goalDiff int
		goalsFor int
	}{
		{"Team D", 15, 15, 25}, // 1st: Highest points, highest goal difference
		{"Team A", 15, 10, 20}, // 2nd: Same points as Team D, but lower goal difference
		{"Team B", 12, 10, 18}, // 3rd: Lower points than Team A and D
		{"Team C", 9, 3, 15},   // 4th: Lowest points
	}

	for i, expected := range expectedOrder {
		assert.Equal(t, expected.name, sortedTeams[i].Name,
			"Team at position %d should be %s", i+1, expected.name)
		assert.Equal(t, expected.points, sortedTeams[i].Stats.Points,
			"Team %s should have %d points", expected.name, expected.points)

		actualGoalDiff := sortedTeams[i].Stats.GoalsFor - sortedTeams[i].Stats.GoalsAgainst
		assert.Equal(t, expected.goalDiff, actualGoalDiff,
			"Team %s should have goal difference of %d", expected.name, expected.goalDiff)
		assert.Equal(t, expected.goalsFor, sortedTeams[i].Stats.GoalsFor,
			"Team %s should have %d goals for", expected.name, expected.goalsFor)
	}

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}
