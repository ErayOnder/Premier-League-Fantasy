package services

import (
	"insider-league/models"
	"insider-league/repository"
)

// TeamService defines the interface for team business logic operations
type TeamService interface {
	Create(team *models.Team) error
	GetAll() ([]models.Team, error)
	GetByID(id int) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id int) error
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

// GetByID retrieves a team by its ID using the repository
func (s *teamService) GetByID(id int) (*models.Team, error) {
	return s.repo.GetByID(id)
}

// Update modifies an existing team using the repository
func (s *teamService) Update(team *models.Team) error {
	return s.repo.Update(team)
}

// Delete removes a team using the repository
func (s *teamService) Delete(id int) error {
	return s.repo.Delete(id)
}
