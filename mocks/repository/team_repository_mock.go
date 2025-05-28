package mocks

import (
	"insider-league/models"
	"insider-league/repository"

	"github.com/stretchr/testify/mock"
)

// MockTeamRepository is a mock implementation of repository.TeamRepository
type MockTeamRepository struct {
	mock.Mock
}

// GetAll mocks the GetAll method
func (m *MockTeamRepository) GetAll() ([]models.Team, error) {
	args := m.Called()
	return args.Get(0).([]models.Team), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockTeamRepository) GetByID(id int) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

// Create mocks the Create method
func (m *MockTeamRepository) Create(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockTeamRepository) Update(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockTeamRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Ensure MockTeamRepository implements repository.TeamRepository
var _ repository.TeamRepository = (*MockTeamRepository)(nil)
