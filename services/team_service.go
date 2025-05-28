package services

import (
	"insider-league/models"
	"insider-league/repository"
	"sort"
)

// TeamService defines the interface for team business logic operations
type TeamService interface {
	Create(team *models.Team) error
	GetAll() ([]models.Team, error)
	GetByID(id int) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id int) error
	GetLeagueTable() ([]models.Team, error)
	UpdateMatchStats(homeTeam, awayTeam *models.Team, homeGoals, awayGoals int, revert bool) error
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

// GetLeagueTable retrieves all teams sorted by points, goal difference, and goals scored
func (s *teamService) GetLeagueTable() ([]models.Team, error) {
	teams, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Sort teams by points (descending), goal difference (descending), and goals scored (descending)
	sort.Slice(teams, func(i, j int) bool {
		// First compare by points
		if teams[i].Stats.Points != teams[j].Stats.Points {
			return teams[i].Stats.Points > teams[j].Stats.Points
		}

		// If points are equal, compare by goal difference
		iGD := teams[i].Stats.GoalsFor - teams[i].Stats.GoalsAgainst
		jGD := teams[j].Stats.GoalsFor - teams[j].Stats.GoalsAgainst
		if iGD != jGD {
			return iGD > jGD
		}

		// If goal difference is equal, compare by goals scored
		return teams[i].Stats.GoalsFor > teams[j].Stats.GoalsFor
	})

	return teams, nil
}

// UpdateMatchStats updates the statistics for both teams based on the match result
// If revert is true, it will subtract the statistics instead of adding them
func (s *teamService) UpdateMatchStats(homeTeam, awayTeam *models.Team, homeGoals, awayGoals int, revert bool) error {
	// Determine the multiplier based on whether we're reverting or applying
	multiplier := 1
	if revert {
		multiplier = -1
	}

	// Update points and match results
	if homeGoals > awayGoals {
		// Home team wins
		homeTeam.Stats.Points += 3 * multiplier
		homeTeam.Stats.Wins += 1 * multiplier
		awayTeam.Stats.Losses += 1 * multiplier
	} else if homeGoals < awayGoals {
		// Away team wins
		awayTeam.Stats.Points += 3 * multiplier
		awayTeam.Stats.Wins += 1 * multiplier
		homeTeam.Stats.Losses += 1 * multiplier
	} else {
		// Draw
		homeTeam.Stats.Points += 1 * multiplier
		awayTeam.Stats.Points += 1 * multiplier
		homeTeam.Stats.Draws += 1 * multiplier
		awayTeam.Stats.Draws += 1 * multiplier
	}

	// Update goal statistics
	homeTeam.Stats.GoalsFor += homeGoals * multiplier
	homeTeam.Stats.GoalsAgainst += awayGoals * multiplier
	awayTeam.Stats.GoalsFor += awayGoals * multiplier
	awayTeam.Stats.GoalsAgainst += homeGoals * multiplier

	// Update teams in database
	if err := s.repo.Update(homeTeam); err != nil {
		return err
	}
	if err := s.repo.Update(awayTeam); err != nil {
		return err
	}

	return nil
}
