package services

import (
	"insider-league/models"
	"insider-league/repository"
)

// MatchService defines the interface for match business logic operations
type MatchService interface {
	Create(match *models.Match) error
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
