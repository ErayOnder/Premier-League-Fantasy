package services

import (
	"insider-league/models"
	"insider-league/repository"
)

// TeamService defines the interface for team business logic operations
type TeamService interface {
	Create(team *models.Team) error
	GetAll() ([]models.Team, error)
}

// teamService implements TeamService interface
type teamService struct {
	repo repository.TeamRepository
}

// NewTeamService creates a new instance of teamService
func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{
		repo: repo,
	}
}

// Create adds a new team using the repository
func (s *teamService) Create(team *models.Team) error {
	return s.repo.Create(team)
}

// GetAll retrieves all teams using the repository
func (s *teamService) GetAll() ([]models.Team, error) {
	return s.repo.GetAll()
}
