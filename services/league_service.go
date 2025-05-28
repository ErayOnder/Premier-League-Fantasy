package services

import (
	"insider-league/helpers"
	"insider-league/models"
	"sort"
)

// LeagueService defines the interface for league-related operations
type LeagueService interface {
	PlayWeek(week int) ([]models.Team, []models.Prediction, error)
	PlayAll() ([]models.Team, error)
	EditMatchResult(matchID int, homeGoals, awayGoals int) (*models.Match, []models.Team, error)
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

// PlayWeek simulates all matches for a given week
func (s *leagueService) PlayWeek(week int) ([]models.Team, []models.Prediction, error) {
	// Get all teams
	teams, err := s.teamService.GetAll()
	if err != nil {
		return nil, nil, err
	}

	// Get all matches for the week
	matches, err := s.matchService.GetByWeek(week)
	if err != nil {
		return nil, nil, err
	}

	// Simulate each match
	for _, match := range matches {
		// Find home and away teams
		var homeTeam, awayTeam *models.Team
		for i := range teams {
			if teams[i].ID == match.HomeTeamID {
				homeTeam = &teams[i]
			}
			if teams[i].ID == match.AwayTeamID {
				awayTeam = &teams[i]
			}
		}

		if homeTeam == nil || awayTeam == nil {
			continue
		}

		// Simulate match score
		homeGoals, awayGoals := helpers.SimulateMatchScore(homeTeam.Strength, awayTeam.Strength)

		// Update match result
		match.HomeTeamScore = homeGoals
		match.AwayTeamScore = awayGoals
		match.IsPlayed = true

		// Update match in database
		if err := s.matchService.Update(&match); err != nil {
			return nil, nil, err
		}

		// Update team statistics
		if err := s.teamService.UpdateMatchStats(homeTeam, awayTeam, homeGoals, awayGoals, false); err != nil {
			return nil, nil, err
		}
	}

	// Get updated league table
	leagueTable, err := s.teamService.GetLeagueTable()
	if err != nil {
		return nil, nil, err
	}

	// Calculate predictions if we're at week 4 or later
	var predictions []models.Prediction
	if week >= 4 {
		predictions = helpers.CalculatePredictions(leagueTable)
	}

	return leagueTable, predictions, nil
}

// PlayAll simulates all remaining matches in the league
func (s *leagueService) PlayAll() ([]models.Team, error) {
	// Get all matches
	matches, err := s.matchService.GetAll()
	if err != nil {
		return nil, err
	}

	// Create a map to track unique weeks
	weekMap := make(map[int]bool)
	for _, match := range matches {
		if !match.IsPlayed {
			weekMap[match.Week] = true
		}
	}

	// Convert map to sorted slice of weeks
	weeks := make([]int, 0, len(weekMap))
	for week := range weekMap {
		weeks = append(weeks, week)
	}
	sort.Ints(weeks)

	// Play each week in order
	for _, week := range weeks {
		_, _, err := s.PlayWeek(week)
		if err != nil {
			return nil, err
		}
	}

	// Get final league table
	leagueTable, err := s.teamService.GetLeagueTable()
	if err != nil {
		return nil, err
	}

	return leagueTable, nil
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
	if err := s.teamService.UpdateMatchStats(homeTeam, awayTeam, match.HomeTeamScore, match.AwayTeamScore, true); err != nil {
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
	if err := s.teamService.UpdateMatchStats(homeTeam, awayTeam, homeGoals, awayGoals, false); err != nil {
		return nil, nil, err
	}

	// Get updated league table
	leagueTable, err := s.teamService.GetLeagueTable()
	if err != nil {
		return nil, nil, err
	}

	return match, leagueTable, nil
}
