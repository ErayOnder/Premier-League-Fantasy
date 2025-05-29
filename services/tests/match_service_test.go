package tests

import (
	repomocks "insider-league/mocks/repository"
	"insider-league/models"
	"insider-league/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchService_Create(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Test data
	newMatch := &models.Match{
		Week:          1,
		HomeTeamID:    1,
		AwayTeamID:    2,
		HomeTeamScore: 0,
		AwayTeamScore: 0,
		IsPlayed:      false,
	}

	// Set up mock expectations
	mockRepo.On("Create", newMatch).Return(nil).Once()

	// Call the function under test
	err := service.Create(newMatch)

	// Assertions
	assert.NoError(t, err, "Create should not return an error")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_GetAll(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Expected matches
	expectedMatches := []models.Match{
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
			HomeTeamScore: 0,
			AwayTeamScore: 3,
			IsPlayed:      true,
		},
	}

	// Set up mock expectations
	mockRepo.On("GetAll").Return(expectedMatches, nil).Once()

	// Call the function under test
	matches, err := service.GetAll()

	// Assertions
	assert.NoError(t, err, "GetAll should not return an error")
	assert.Equal(t, expectedMatches, matches, "Matches should match expected")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_GetByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Test data
	matchID := 1
	expectedMatch := &models.Match{
		ID:            1,
		Week:          1,
		HomeTeamID:    1,
		AwayTeamID:    2,
		HomeTeamScore: 2,
		AwayTeamScore: 1,
		IsPlayed:      true,
	}

	// Set up mock expectations
	mockRepo.On("GetByID", matchID).Return(expectedMatch, nil).Once()

	// Call the function under test
	match, err := service.GetByID(matchID)

	// Assertions
	assert.NoError(t, err, "GetByID should not return an error")
	assert.Equal(t, expectedMatch, match, "Match should match expected")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_GetByWeek(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Test data
	week := 2
	expectedMatches := []models.Match{
		{
			ID:            3,
			Week:          2,
			HomeTeamID:    1,
			AwayTeamID:    3,
			HomeTeamScore: 1,
			AwayTeamScore: 0,
			IsPlayed:      true,
		},
		{
			ID:            4,
			Week:          2,
			HomeTeamID:    2,
			AwayTeamID:    4,
			HomeTeamScore: 2,
			AwayTeamScore: 2,
			IsPlayed:      true,
		},
	}

	// Set up mock expectations
	mockRepo.On("GetByWeek", week).Return(expectedMatches, nil).Once()

	// Call the function under test
	matches, err := service.GetByWeek(week)

	// Assertions
	assert.NoError(t, err, "GetByWeek should not return an error")
	assert.Equal(t, expectedMatches, matches, "Matches should match expected")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_GetUnplayedWeeks(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Expected unplayed weeks
	expectedWeeks := []int{3, 4, 5}

	// Set up mock expectations
	mockRepo.On("GetUnplayedWeeks").Return(expectedWeeks, nil).Once()

	// Call the function under test
	weeks, err := service.GetUnplayedWeeks()

	// Assertions
	assert.NoError(t, err, "GetUnplayedWeeks should not return an error")
	assert.Equal(t, expectedWeeks, weeks, "Weeks should match expected")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_Update(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Test data
	updatedMatch := &models.Match{
		ID:            1,
		Week:          1,
		HomeTeamID:    1,
		AwayTeamID:    2,
		HomeTeamScore: 3,
		AwayTeamScore: 1,
		IsPlayed:      true,
	}

	// Set up mock expectations
	mockRepo.On("Update", updatedMatch).Return(nil).Once()

	// Call the function under test
	err := service.Update(updatedMatch)

	// Assertions
	assert.NoError(t, err, "Update should not return an error")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}

func TestMatchService_Delete(t *testing.T) {
	// Create mock repository
	mockRepo := new(repomocks.MockMatchRepository)

	// Create match service with mock
	service := services.NewMatchService(mockRepo)

	// Test data
	matchID := 1

	// Set up mock expectations
	mockRepo.On("Delete", matchID).Return(nil).Once()

	// Call the function under test
	err := service.Delete(matchID)

	// Assertions
	assert.NoError(t, err, "Delete should not return an error")

	// Verify that all expected calls were made
	mockRepo.AssertExpectations(t)
}
