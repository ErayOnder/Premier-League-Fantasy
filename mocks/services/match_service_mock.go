package mocks

import (
	"insider-league/models"

	"github.com/stretchr/testify/mock"
)

// MockMatchService is a mock implementation of MatchService interface
type MockMatchService struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockMatchService) Create(match *models.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

// GetAll mocks the GetAll method
func (m *MockMatchService) GetAll() ([]models.Match, error) {
	args := m.Called()
	return args.Get(0).([]models.Match), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockMatchService) GetByID(id int) (*models.Match, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

// GetByWeek mocks the GetByWeek method
func (m *MockMatchService) GetByWeek(week int) ([]models.Match, error) {
	args := m.Called(week)
	return args.Get(0).([]models.Match), args.Error(1)
}

// Update mocks the Update method
func (m *MockMatchService) Update(match *models.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockMatchService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetUnplayedWeeks mocks the GetUnplayedWeeks method
func (m *MockMatchService) GetUnplayedWeeks() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}
