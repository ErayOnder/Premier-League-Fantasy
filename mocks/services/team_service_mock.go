package mocks

import (
	"insider-league/models"

	"github.com/stretchr/testify/mock"
)

// MockTeamService is a mock of TeamService interface
type MockTeamService struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockTeamService) Create(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

// GetAll mocks the GetAll method
func (m *MockTeamService) GetAll() ([]models.Team, error) {
	args := m.Called()
	return args.Get(0).([]models.Team), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockTeamService) GetByID(id int) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

// Update mocks the Update method
func (m *MockTeamService) Update(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockTeamService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetLeagueTable mocks the GetLeagueTable method
func (m *MockTeamService) GetLeagueTable() ([]models.Team, error) {
	args := m.Called()
	return args.Get(0).([]models.Team), args.Error(1)
}

// UpdateTeamStats mocks the UpdateTeamStats method
func (m *MockTeamService) UpdateTeamStats(homeTeam, awayTeam *models.Team, homeGoals, awayGoals int, revert bool) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, revert)
	return args.Error(0)
}

// GetTeamRankings mocks base method
func (m *MockTeamService) GetTeamRankings() ([]models.Team, error) {
	args := m.Called()
	return args.Get(0).([]models.Team), args.Error(1)
}
