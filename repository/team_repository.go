package repository

import (
	"insider-league/models"

	"gorm.io/gorm"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	GetAll() ([]models.Team, error)
	GetByID(id int) (*models.Team, error)
	Create(team *models.Team) error
	Update(team *models.Team) error
	Delete(id int) error
}

// teamRepository implements TeamRepository interface
type teamRepository struct {
	db *gorm.DB
}

// NewTeamRepository creates a new instance of teamRepository
func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		db: db,
	}
}

// GetAll retrieves all teams from the database
func (r *teamRepository) GetAll() ([]models.Team, error) {
	var teams []models.Team
	result := r.db.Find(&teams)
	return teams, result.Error
}

// GetByID retrieves a team by its ID
func (r *teamRepository) GetByID(id int) (*models.Team, error) {
	var team models.Team
	result := r.db.First(&team, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

// Create adds a new team to the database
func (r *teamRepository) Create(team *models.Team) error {
	result := r.db.Create(team)
	return result.Error
}

// Update modifies an existing team in the database
func (r *teamRepository) Update(team *models.Team) error {
	result := r.db.Save(team)
	return result.Error
}

// Delete removes a team from the database by its ID
func (r *teamRepository) Delete(id int) error {
	result := r.db.Delete(&models.Team{}, id)
	return result.Error
}
