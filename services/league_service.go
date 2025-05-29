package services

import (
	"insider-league/helpers"
	"insider-league/models"
)

// LeagueService defines the interface for league-related operations
type LeagueService interface {
	GetLeagueTable() ([]models.Team, error)
	PlayWeeks(playAll bool) ([]models.Team, []models.Match, []models.Prediction, error)
	GetWeekResults(week int) ([]models.Match, error)
	EditMatchResult(matchID int, homeGoals, awayGoals int) (*models.Match, []models.Team, error)
	ResetLeague() error
}

// leagueService implements the LeagueService interface
type leagueService struct {
	teamService  TeamService
	matchService MatchService
}

// NewLeagueService creates a new instance of leagueService
func NewLeagueService(teamService TeamService, matchService MatchService) LeagueService {
	return &leagueService{
		teamService:  teamService,
		matchService: matchService,
	}
}

// GetLeagueTable retrieves the current league table
func (s *leagueService) GetLeagueTable() ([]models.Team, error) {
	return s.teamService.GetTeamRankings()
}

// PlayWeeks simulates weeks based on the playAll parameter
// If playAll is false, it plays only the next unplayed week
// If playAll is true, it plays all remaining unplayed weeks
func (s *leagueService) PlayWeeks(playAll bool) ([]models.Team, []models.Match, []models.Prediction, error) {
	// Get all unplayed weeks sorted
	unplayedWeeks, err := s.matchService.GetUnplayedWeeks()
	if err != nil {
		return nil, nil, nil, err
	}

	// If no unplayed weeks found, return current league table
	if len(unplayedWeeks) == 0 {
		leagueTable, err := s.teamService.GetTeamRankings()
		if err != nil {
			return nil, nil, nil, err
		}
		return leagueTable, []models.Match{}, []models.Prediction{}, nil
	}

	var allMatches []models.Match
	var currentWeek int

	// Loop through unplayed weeks
	for _, week := range unplayedWeeks {
		currentWeek = week

		// Get matches for the current week
		weekMatches, err := s.matchService.GetByWeek(week)
		if err != nil {
			return nil, nil, nil, err
		}

		// Simulate each match in the week
		for i := range weekMatches {
			match := &weekMatches[i]

			// Use the preloaded teams
			homeTeam := &match.HomeTeam
			awayTeam := &match.AwayTeam

			// Simulate match score
			homeGoals, awayGoals := helpers.SimulateMatchScore(homeTeam.Strength, awayTeam.Strength)

			// Update match result
			match.HomeTeamScore = homeGoals
			match.AwayTeamScore = awayGoals
			match.IsPlayed = true

			// Update match in database
			if err := s.matchService.Update(match); err != nil {
				return nil, nil, nil, err
			}

			// Update team statistics
			if err := s.teamService.UpdateTeamStats(homeTeam, awayTeam, homeGoals, awayGoals, false); err != nil {
				return nil, nil, nil, err
			}
		}

		// Add week matches to all matches
		allMatches = append(allMatches, weekMatches...)

		// If not playing all weeks, break after first week
		if !playAll {
			break
		}
	}

	// Get updated league table
	leagueTable, err := s.teamService.GetTeamRankings()
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate predictions if we're at week 4 or later
	var predictions []models.Prediction
	if currentWeek >= 4 {
		predictions = helpers.CalculatePredictions(leagueTable, currentWeek)
	}

	return leagueTable, allMatches, predictions, nil
}

// GetWeekResults retrieves the results for a specific week
func (s *leagueService) GetWeekResults(week int) ([]models.Match, error) {
	matches, err := s.matchService.GetByWeek(week)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// EditMatchResult updates a match result and recalculates team statistics
func (s *leagueService) EditMatchResult(matchID int, homeGoals, awayGoals int) (*models.Match, []models.Team, error) {
	// Get match with preloaded teams
	match, err := s.matchService.GetByID(matchID)
	if err != nil {
		return nil, nil, err
	}

	// Get teams
	homeTeam, err := s.teamService.GetByID(int(match.HomeTeamID))
	if err != nil {
		return nil, nil, err
	}
	awayTeam, err := s.teamService.GetByID(int(match.AwayTeamID))
	if err != nil {
		return nil, nil, err
	}

	// Revert Phase: Undo the effects of the original match result
	if err := s.teamService.UpdateTeamStats(homeTeam, awayTeam, match.HomeTeamScore, match.AwayTeamScore, true); err != nil {
		return nil, nil, err
	}

	// Apply Phase: Apply the new match result
	match.HomeTeamScore = homeGoals
	match.AwayTeamScore = awayGoals
	match.IsPlayed = true

	// Update match in database
	if err := s.matchService.Update(match); err != nil {
		return nil, nil, err
	}

	// Update team statistics with new result
	if err := s.teamService.UpdateTeamStats(homeTeam, awayTeam, homeGoals, awayGoals, false); err != nil {
		return nil, nil, err
	}

	// Get updated league table
	leagueTable, err := s.teamService.GetTeamRankings()
	if err != nil {
		return nil, nil, err
	}

	return match, leagueTable, nil
}

// ResetLeague resets all match results and team statistics
func (s *leagueService) ResetLeague() error {
	// Reset all matches
	matches, err := s.matchService.GetAll()
	if err != nil {
		return err
	}

	for _, match := range matches {
		match.HomeTeamScore = 0
		match.AwayTeamScore = 0
		match.IsPlayed = false
		if err := s.matchService.Update(&match); err != nil {
			return err
		}
	}

	// Reset all teams
	teams, err := s.teamService.GetAll()
	if err != nil {
		return err
	}

	for _, team := range teams {
		team.Stats = models.Stats{
			Points:         0,
			GoalsFor:       0,
			GoalsAgainst:   0,
			GoalDifference: 0,
			Wins:           0,
			Draws:          0,
			Losses:         0,
		}
		if err := s.teamService.Update(&team); err != nil {
			return err
		}
	}

	return nil
}
