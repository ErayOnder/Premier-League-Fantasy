package services

import (
	"insider-league/helpers"
	"insider-league/models"
	"insider-league/repository"
)

// LeagueService defines the interface for league-related operations
type LeagueService interface {
	PlayWeek(week int) ([]models.Team, []models.Prediction, error)
}

// leagueService implements the LeagueService interface
type leagueService struct {
	teamRepo  repository.TeamRepository
	matchRepo repository.MatchRepository
}

// NewLeagueService creates a new instance of leagueService
func NewLeagueService(teamRepo repository.TeamRepository, matchRepo repository.MatchRepository) LeagueService {
	return &leagueService{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
	}
}

// PlayWeek simulates all matches for a given week
func (s *leagueService) PlayWeek(week int) ([]models.Team, []models.Prediction, error) {
	// Get all teams
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, nil, err
	}

	// Get all matches for the week
	matches, err := s.matchRepo.GetByWeek(week)
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
		if err := s.matchRepo.Update(&match); err != nil {
			return nil, nil, err
		}

		// Update team statistics
		if homeGoals > awayGoals {
			homeTeam.Stats.Points += 3
			homeTeam.Stats.Wins++
			awayTeam.Stats.Losses++
		} else if homeGoals < awayGoals {
			awayTeam.Stats.Points += 3
			awayTeam.Stats.Wins++
			homeTeam.Stats.Losses++
		} else {
			homeTeam.Stats.Points++
			awayTeam.Stats.Points++
			homeTeam.Stats.Draws++
			awayTeam.Stats.Draws++
		}

		homeTeam.Stats.GoalsFor += homeGoals
		homeTeam.Stats.GoalsAgainst += awayGoals
		awayTeam.Stats.GoalsFor += awayGoals
		awayTeam.Stats.GoalsAgainst += homeGoals

		// Update teams in database
		if err := s.teamRepo.Update(homeTeam); err != nil {
			return nil, nil, err
		}
		if err := s.teamRepo.Update(awayTeam); err != nil {
			return nil, nil, err
		}
	}

	// Get updated league table
	leagueTable, err := s.teamRepo.GetAll()
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
