package services

import (
	"insider-league/models"
	"insider-league/repository"
)

// MatchService defines the interface for match business logic operations
type MatchService interface {
	Create(match *models.Match) error
	GetAll() ([]models.Match, error)
	GetByID(id int) (*models.Match, error)
	GetByWeek(week int) ([]models.Match, error)
	GetUnplayedWeeks() ([]int, error)
	Update(match *models.Match) error
	Delete(id int) error
}

// matchService implements MatchService interface
type matchService struct {
	repo repository.MatchRepository
}

// NewMatchService creates a new instance of matchService
func NewMatchService(repo repository.MatchRepository) MatchService {
	return &matchService{
		repo: repo,
	}
}

// Create adds a new match using the repository
func (s *matchService) Create(match *models.Match) error {
	return s.repo.Create(match)
}

// GetAll retrieves all matches using the repository
func (s *matchService) GetAll() ([]models.Match, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a match by its ID using the repository
func (s *matchService) GetByID(id int) (*models.Match, error) {
	return s.repo.GetByID(id)
}

// GetByWeek retrieves all matches for a specific week using the repository
func (s *matchService) GetByWeek(week int) ([]models.Match, error) {
	return s.repo.GetByWeek(week)
}

// GetUnplayedWeeks retrieves all unplayed weeks sorted
func (s *matchService) GetUnplayedWeeks() ([]int, error) {
	return s.repo.GetUnplayedWeeks()
}

// Update modifies an existing match using the repository
func (s *matchService) Update(match *models.Match) error {
	return s.repo.Update(match)
}

// Delete removes a match using the repository
func (s *matchService) Delete(id int) error {
	return s.repo.Delete(id)
}
