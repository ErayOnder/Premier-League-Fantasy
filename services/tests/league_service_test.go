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
	mockTeamService.On("UpdateTeamStats", homeTeam, awayTeam, 2, 0, true).Return(nil).Once()

	// Update match expectation
	mockMatchService.On("Update", mock.MatchedBy(func(match *models.Match) bool {
		return match.ID == 1 && match.HomeTeamScore == newHomeGoals && match.AwayTeamScore == newAwayGoals && match.IsPlayed == true
	})).Return(nil).Once()

	// Second call: Apply the new match result (3-1) with revert=false
	mockTeamService.On("UpdateTeamStats", homeTeam, awayTeam, newHomeGoals, newAwayGoals, false).Return(nil).Once()

	// Get league table expectation
	mockTeamService.On("GetTeamRankings").Return(expectedLeagueTable, nil).Once()

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

func TestLeagueService_PlayWeeks_NextWeek(t *testing.T) {
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
			homeTeam := models.Team{
				ID:       1,
				Name:     "Team A",
				Strength: 80,
				Stats: models.Stats{
					Points:         10,
					Wins:           3,
					Draws:          1,
					Losses:         0,
					GoalsFor:       8,
					GoalsAgainst:   2,
					GoalDifference: 6,
				},
			}

			awayTeam := models.Team{
				ID:       2,
				Name:     "Team B",
				Strength: 75,
				Stats: models.Stats{
					Points:         7,
					Wins:           2,
					Draws:          1,
					Losses:         1,
					GoalsFor:       6,
					GoalsAgainst:   4,
					GoalDifference: 2,
				},
			}

			// Create sample unplayed matches for the week with preloaded teams
			matches := []models.Match{
				{
					ID:            1,
					Week:          tt.week,
					HomeTeamID:    1,
					AwayTeamID:    2,
					HomeTeamScore: 0,
					AwayTeamScore: 0,
					IsPlayed:      false,
					HomeTeam:      homeTeam,
					AwayTeam:      awayTeam,
				},
				{
					ID:            2,
					Week:          tt.week,
					HomeTeamID:    2,
					AwayTeamID:    1,
					HomeTeamScore: 0,
					AwayTeamScore: 0,
					IsPlayed:      false,
					HomeTeam:      awayTeam,
					AwayTeam:      homeTeam,
				},
			}

			// Expected league table after playing the week
			expectedLeagueTable := []models.Team{homeTeam, awayTeam}

			// Set up mock expectations
			mockMatchService.On("GetUnplayedWeeks").Return([]int{tt.week}, nil).Once()
			mockMatchService.On("GetByWeek", tt.week).Return(matches, nil).Once()

			// For each match, expect Update to be called
			for i := range matches {
				mockMatchService.On("Update", mock.MatchedBy(func(match *models.Match) bool {
					return match.ID == matches[i].ID && match.IsPlayed == true
				})).Return(nil).Once()
			}

			// For each match, expect UpdateTeamStats to be called
			for range matches {
				mockTeamService.On("UpdateTeamStats",
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

			// Expect GetTeamRankings to be called
			mockTeamService.On("GetTeamRankings").Return(expectedLeagueTable, nil).Once()

			// Call the function under test - play only next week
			leagueTable, returnedMatches, predictions, err := service.PlayWeeks(false)

			// Assertions
			assert.NoError(t, err, "PlayWeeks should not return an error")
			assert.Equal(t, expectedLeagueTable, leagueTable, "League table should match expected")
			assert.NotNil(t, returnedMatches, "Returned matches should not be nil")
			assert.Len(t, returnedMatches, len(matches), "Should return all matches for the week")

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

func TestLeagueService_PlayWeeks_NoUnplayedWeeks(t *testing.T) {
	// Create mock services
	mockTeamService := new(servicemocks.MockTeamService)
	mockMatchService := new(servicemocks.MockMatchService)

	// Create league service with mocks
	service := services.NewLeagueService(mockTeamService, mockMatchService)

	// Expected league table when no unplayed weeks remain
	expectedLeagueTable := []models.Team{
		{
			ID:   1,
			Name: "Team A",
			Stats: models.Stats{
				Points:       15,
				Wins:         5,
				Draws:        0,
				Losses:       0,
				GoalsFor:     12,
				GoalsAgainst: 3,
			},
		},
		{
			ID:   2,
			Name: "Team B",
			Stats: models.Stats{
				Points:       12,
				Wins:         4,
				Draws:        0,
				Losses:       1,
				GoalsFor:     10,
				GoalsAgainst: 5,
			},
		},
	}

	// Set up mock expectations - return empty slice for no unplayed weeks
	mockMatchService.On("GetUnplayedWeeks").Return([]int{}, nil).Once()
	mockTeamService.On("GetTeamRankings").Return(expectedLeagueTable, nil).Once()

	// Call the function under test
	leagueTable, returnedMatches, predictions, err := service.PlayWeeks(false)

	// Assertions - should NOT return an error, but return current state
	assert.NoError(t, err, "PlayWeeks should not return an error when no unplayed weeks remain")
	assert.Equal(t, expectedLeagueTable, leagueTable, "League table should match expected")
	assert.NotNil(t, returnedMatches, "Returned matches should not be nil")
	assert.Empty(t, returnedMatches, "Returned matches should be empty when no unplayed weeks")
	assert.NotNil(t, predictions, "Predictions should not be nil")
	assert.Empty(t, predictions, "Predictions should be empty when no unplayed weeks")

	// Verify that the expected calls were made
	mockMatchService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}

func TestLeagueService_GetLeagueTable(t *testing.T) {
	// Create mock services
	mockTeamService := new(servicemocks.MockTeamService)
	mockMatchService := new(servicemocks.MockMatchService)

	// Create league service with mocks
	service := services.NewLeagueService(mockTeamService, mockMatchService)

	// Expected league table
	expectedLeagueTable := []models.Team{
		{
			ID:   1,
			Name: "Team A",
			Stats: models.Stats{
				Points:       15,
				Wins:         5,
				Draws:        0,
				Losses:       0,
				GoalsFor:     12,
				GoalsAgainst: 3,
			},
		},
		{
			ID:   2,
			Name: "Team B",
			Stats: models.Stats{
				Points:       12,
				Wins:         4,
				Draws:        0,
				Losses:       1,
				GoalsFor:     10,
				GoalsAgainst: 5,
			},
		},
	}

	// Set up mock expectations
	mockTeamService.On("GetTeamRankings").Return(expectedLeagueTable, nil).Once()

	// Call the function under test
	leagueTable, err := service.GetLeagueTable()

	// Assertions
	assert.NoError(t, err, "GetLeagueTable should not return an error")
	assert.Equal(t, expectedLeagueTable, leagueTable, "League table should match expected")

	// Verify that the expected calls were made
	mockTeamService.AssertExpectations(t)
	mockMatchService.AssertExpectations(t)
}

func TestLeagueService_GetWeekResults(t *testing.T) {
	// Create mock services
	mockTeamService := new(servicemocks.MockTeamService)
	mockMatchService := new(servicemocks.MockMatchService)

	// Create league service with mocks
	service := services.NewLeagueService(mockTeamService, mockMatchService)

	// Test data
	week := 3
	expectedMatches := []models.Match{
		{
			ID:            1,
			Week:          3,
			HomeTeamID:    1,
			AwayTeamID:    2,
			HomeTeamScore: 2,
			AwayTeamScore: 1,
			IsPlayed:      true,
		},
		{
			ID:            2,
			Week:          3,
			HomeTeamID:    3,
			AwayTeamID:    4,
			HomeTeamScore: 0,
			AwayTeamScore: 3,
			IsPlayed:      true,
		},
	}

	// Set up mock expectations
	mockMatchService.On("GetByWeek", week).Return(expectedMatches, nil).Once()

	// Call the function under test
	matches, err := service.GetWeekResults(week)

	// Assertions
	assert.NoError(t, err, "GetWeekResults should not return an error")
	assert.Equal(t, expectedMatches, matches, "Matches should match expected")

	// Verify that the expected calls were made
	mockMatchService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}

func TestLeagueService_ResetLeague(t *testing.T) {
	// Create mock services
	mockTeamService := new(servicemocks.MockTeamService)
	mockMatchService := new(servicemocks.MockMatchService)

	// Create league service with mocks
	service := services.NewLeagueService(mockTeamService, mockMatchService)

	// Mock data - existing matches with played results
	existingMatches := []models.Match{
		{
			ID:            1,
			Week:          1,
			HomeTeamID:    1,
			AwayTeamID:    2,
			HomeTeamScore: 2,
			AwayTeamScore: 1,
			IsPlayed:      true,
		},
		{
			ID:            2,
			Week:          1,
			HomeTeamID:    3,
			AwayTeamID:    4,
			HomeTeamScore: 1,
			AwayTeamScore: 0,
			IsPlayed:      true,
		},
	}

	// Mock data - existing teams with stats
	existingTeams := []models.Team{
		{
			ID:   1,
			Name: "Team A",
			Stats: models.Stats{
				Points:       15,
				Wins:         5,
				Draws:        0,
				Losses:       0,
				GoalsFor:     12,
				GoalsAgainst: 3,
			},
		},
		{
			ID:   2,
			Name: "Team B",
			Stats: models.Stats{
				Points:       12,
				Wins:         4,
				Draws:        0,
				Losses:       1,
				GoalsFor:     10,
				GoalsAgainst: 5,
			},
		},
	}

	// Set up mock expectations
	mockMatchService.On("GetAll").Return(existingMatches, nil).Once()
	mockTeamService.On("GetAll").Return(existingTeams, nil).Once()

	// Expect Update to be called for each match (reset to unplayed state)
	for _, match := range existingMatches {
		mockMatchService.On("Update", mock.MatchedBy(func(m *models.Match) bool {
			return m.ID == match.ID && m.HomeTeamScore == 0 && m.AwayTeamScore == 0 && m.IsPlayed == false
		})).Return(nil).Once()
	}

	// Expect Update to be called for each team (reset stats)
	for _, team := range existingTeams {
		mockTeamService.On("Update", mock.MatchedBy(func(t *models.Team) bool {
			return t.ID == team.ID &&
				t.Stats.Points == 0 && t.Stats.Wins == 0 && t.Stats.Draws == 0 &&
				t.Stats.Losses == 0 && t.Stats.GoalsFor == 0 && t.Stats.GoalsAgainst == 0 &&
				t.Stats.GoalDifference == 0
		})).Return(nil).Once()
	}

	// Call the function under test
	err := service.ResetLeague()

	// Assertions
	assert.NoError(t, err, "ResetLeague should not return an error")

	// Verify that all expected calls were made
	mockMatchService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}
