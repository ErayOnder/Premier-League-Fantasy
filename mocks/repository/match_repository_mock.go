package mocks

import (
	"insider-league/models"
	"insider-league/repository"

	"github.com/stretchr/testify/mock"
)

// MockMatchRepository is a mock implementation of repository.MatchRepository
type MockMatchRepository struct {
	mock.Mock
}

// GetAll mocks the GetAll method
func (m *MockMatchRepository) GetAll() ([]models.Match, error) {
	args := m.Called()
	return args.Get(0).([]models.Match), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockMatchRepository) GetByID(id int) (*models.Match, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

// GetByWeek mocks the GetByWeek method
func (m *MockMatchRepository) GetByWeek(week int) ([]models.Match, error) {
	args := m.Called(week)
	return args.Get(0).([]models.Match), args.Error(1)
}

// Create mocks the Create method
func (m *MockMatchRepository) Create(match *models.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockMatchRepository) Update(match *models.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockMatchRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Ensure MockMatchRepository implements repository.MatchRepository
var _ repository.MatchRepository = (*MockMatchRepository)(nil)
