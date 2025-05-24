package repository

import (
	"insider-league/models"

	"gorm.io/gorm"
)

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	GetAll() ([]models.Match, error)
	GetByID(id int) (*models.Match, error)
	GetByWeek(week int) ([]models.Match, error)
	Create(match *models.Match) error
	Update(match *models.Match) error
	Delete(id int) error
}

// matchRepository implements MatchRepository interface
type matchRepository struct {
	db *gorm.DB
}

// NewMatchRepository creates a new instance of matchRepository
func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{
		db: db,
	}
}

// GetAll retrieves all matches from the database with their associated teams
func (r *matchRepository) GetAll() ([]models.Match, error) {
	var matches []models.Match
	result := r.db.Preload("HomeTeam").Preload("AwayTeam").Find(&matches)
	return matches, result.Error
}

// GetByID retrieves a match by its ID with its associated teams
func (r *matchRepository) GetByID(id int) (*models.Match, error) {
	var match models.Match
	result := r.db.Preload("HomeTeam").Preload("AwayTeam").First(&match, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &match, nil
}

// GetByWeek retrieves all matches for a specific week with their associated teams
func (r *matchRepository) GetByWeek(week int) ([]models.Match, error) {
	var matches []models.Match
	result := r.db.Preload("HomeTeam").Preload("AwayTeam").Where("week = ?", week).Find(&matches)
	return matches, result.Error
}

// Create adds a new match to the database
func (r *matchRepository) Create(match *models.Match) error {
	result := r.db.Create(match)
	return result.Error
}

// Update modifies an existing match in the database
func (r *matchRepository) Update(match *models.Match) error {
	result := r.db.Save(match)
	return result.Error
}

// Delete removes a match from the database by its ID
func (r *matchRepository) Delete(id int) error {
	result := r.db.Delete(&models.Match{}, id)
	return result.Error
}
