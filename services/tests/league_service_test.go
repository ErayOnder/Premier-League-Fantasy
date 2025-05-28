package tests

import (
	servicemocks "insider-league/mocks/services"
	"insider-league/models"
	"insider-league/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeagueService_EditMatchResult(t *testing.T) {
	// Create mock services
	mockTeamService := new(servicemocks.MockTeamService)
	mockMatchService := new(servicemocks.MockMatchService)

	// Create league service with mocks
	service := services.NewLeagueService(mockTeamService, mockMatchService)

	// Test data
	matchID := 1
	newHomeGoals := 3
	newAwayGoals := 1

	// Create sample match (already played with original score)
	originalMatch := &models.Match{
		ID:            1,
		Week:          1,
		HomeTeamID:    1,
		AwayTeamID:    2,
		HomeTeamScore: 2, // Original score
		AwayTeamScore: 0, // Original score
		IsPlayed:      true,
	}

	// Create sample teams
	homeTeam := &models.Team{
		ID:   1,
		Name: "Home Team",
		Stats: models.Stats{
			Points:       13, // Some existing stats
			Wins:         3,
			Draws:        1,
			Losses:       1,
			GoalsFor:     8,
			GoalsAgainst: 4,
		},
	}

	awayTeam := &models.Team{
		ID:   2,
		Name: "Away Team",
		Stats: models.Stats{
			Points:       6,
			Wins:         2,
			Draws:        0,
			Losses:       3,
			GoalsFor:     5,
			GoalsAgainst: 8,
		},
	}

	// Create expected league table
	expectedLeagueTable := []models.Team{*homeTeam, *awayTeam}

	// Set up mock expectations
	mockMatchService.On("GetByID", matchID).Return(originalMatch, nil).Once()
	mockTeamService.On("GetByID", 1).Return(homeTeam, nil).Once()
	mockTeamService.On("GetByID", 2).Return(awayTeam, nil).Once()

	// First call: Revert the original match result (2-0) with revert=true
	mockTeamService.On("UpdateMatchStats", homeTeam, awayTeam, 2, 0, true).Return(nil).Once()

	// Update match expectation
	mockMatchService.On("Update", mock.MatchedBy(func(match *models.Match) bool {
		return match.ID == 1 && match.HomeTeamScore == newHomeGoals && match.AwayTeamScore == newAwayGoals && match.IsPlayed == true
	})).Return(nil).Once()

	// Second call: Apply the new match result (3-1) with revert=false
	mockTeamService.On("UpdateMatchStats", homeTeam, awayTeam, newHomeGoals, newAwayGoals, false).Return(nil).Once()

	// Get league table expectation
	mockTeamService.On("GetLeagueTable").Return(expectedLeagueTable, nil).Once()

	// Call the function under test
	updatedMatch, leagueTable, err := service.EditMatchResult(matchID, newHomeGoals, newAwayGoals)

	// Assertions
	assert.NoError(t, err, "EditMatchResult should not return an error")
	assert.NotNil(t, updatedMatch, "Updated match should not be nil")
	assert.Equal(t, newHomeGoals, updatedMatch.HomeTeamScore, "Home team score should be updated")
	assert.Equal(t, newAwayGoals, updatedMatch.AwayTeamScore, "Away team score should be updated")
	assert.True(t, updatedMatch.IsPlayed, "Match should be marked as played")
	assert.Equal(t, expectedLeagueTable, leagueTable, "League table should match expected")

	// Verify that all expected calls were made in the correct order
	mockMatchService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}

func TestLeagueService_PlayWeek(t *testing.T) {
	tests := []struct {
		name                   string
		week                   int
		expectPredictions      bool
		expectedPredictionsLen int
		description            string
	}{
		{
			name:                   "Week 3 - No predictions",
			week:                   3,
			expectPredictions:      false,
			expectedPredictionsLen: 0,
			description:            "Week < 4 should not generate predictions",
		},
		{
			name:                   "Week 4 - With predictions",
			week:                   4,
			expectPredictions:      true,
			expectedPredictionsLen: 1, // We'll mock this to return 1 prediction
			description:            "Week >= 4 should generate predictions",
		},
		{
			name:                   "Week 10 - With predictions",
			week:                   10,
			expectPredictions:      true,
			expectedPredictionsLen: 1,
			description:            "Week >= 4 should generate predictions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			mockTeamService := new(servicemocks.MockTeamService)
			mockMatchService := new(servicemocks.MockMatchService)

			// Create league service with mocks
			service := services.NewLeagueService(mockTeamService, mockMatchService)

			// Create sample teams
			teams := []models.Team{
				{
					ID:       1,
					Name:     "Team A",
					Strength: 80,
					Stats: models.Stats{
						Points:       10,
						Wins:         3,
						Draws:        1,
						Losses:       0,
						GoalsFor:     8,
						GoalsAgainst: 2,
					},
				},
				{
					ID:       2,
					Name:     "Team B",
					Strength: 75,
					Stats: models.Stats{
						Points:       7,
						Wins:         2,
						Draws:        1,
						Losses:       1,
						GoalsFor:     6,
						GoalsAgainst: 4,
					},
				},
			}

			// Create sample unplayed matches for the week
			matches := []models.Match{
				{
					ID:            1,
					Week:          tt.week,
					HomeTeamID:    1,
					AwayTeamID:    2,
					HomeTeamScore: 0,
					AwayTeamScore: 0,
					IsPlayed:      false,
				},
				{
					ID:            2,
					Week:          tt.week,
					HomeTeamID:    2,
					AwayTeamID:    1,
					HomeTeamScore: 0,
					AwayTeamScore: 0,
					IsPlayed:      false,
				},
			}

			// Expected league table after playing the week
			expectedLeagueTable := teams // Simplified for test

			// Set up mock expectations
			mockTeamService.On("GetAll").Return(teams, nil).Once()
			mockMatchService.On("GetByWeek", tt.week).Return(matches, nil).Once()

			// For each match, expect Update to be called
			for i := range matches {
				mockMatchService.On("Update", mock.MatchedBy(func(match *models.Match) bool {
					return match.ID == matches[i].ID && match.IsPlayed == true
				})).Return(nil).Once()
			}

			// For each match, expect UpdateMatchStats to be called
			// We need to match the teams by ID since the function finds them in the loop
			for range matches {
				mockTeamService.On("UpdateMatchStats",
					mock.MatchedBy(func(homeTeam *models.Team) bool {
						return homeTeam.ID == 1 || homeTeam.ID == 2
					}),
					mock.MatchedBy(func(awayTeam *models.Team) bool {
						return awayTeam.ID == 1 || awayTeam.ID == 2
					}),
					mock.AnythingOfType("int"), // homeGoals (simulated)
					mock.AnythingOfType("int"), // awayGoals (simulated)
					false,                      // revert = false
				).Return(nil).Once()
			}

			// Expect GetLeagueTable to be called
			mockTeamService.On("GetLeagueTable").Return(expectedLeagueTable, nil).Once()

			// Call the function under test
			leagueTable, predictions, err := service.PlayWeek(tt.week)

			// Assertions
			assert.NoError(t, err, "PlayWeek should not return an error")
			assert.Equal(t, expectedLeagueTable, leagueTable, "League table should match expected")

			if tt.expectPredictions {
				assert.NotEmpty(t, predictions, "Predictions should not be empty for week >= 4")
				// Note: The actual predictions are generated by helpers.CalculatePredictions
				// which we're not mocking, so we just verify they exist
			} else {
				assert.Empty(t, predictions, "Predictions should be empty for week < 4")
			}

			// Verify that all expected calls were made
			mockMatchService.AssertExpectations(t)
			mockTeamService.AssertExpectations(t)
		})
	}
}
