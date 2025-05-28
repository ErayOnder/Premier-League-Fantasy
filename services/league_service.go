package services

import (
	"insider-league/helpers"
	"insider-league/models"
	"insider-league/repository"
	"sort"
)

// LeagueService defines the interface for league operations
type LeagueService interface {
	PlayWeek(week int) ([]models.Team, error)
}

// leagueService implements LeagueService interface
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

// PlayWeek simulates all matches for a given week and returns the updated league table
func (s *leagueService) PlayWeek(week int) ([]models.Team, error) {
	// Get all matches for the specified week
	matches, err := s.matchRepo.GetByWeek(week)
	if err != nil {
		return nil, err
	}

	// Process each match
	for _, match := range matches {
		if !match.IsPlayed {
			// Get full team objects
			homeTeam, err := s.teamRepo.GetByID(int(match.HomeTeamID))
			if err != nil {
				return nil, err
			}

			awayTeam, err := s.teamRepo.GetByID(int(match.AwayTeamID))
			if err != nil {
				return nil, err
			}

			// Simulate match score
			homeGoals, awayGoals := helpers.SimulateMatchScore(homeTeam.Strength, awayTeam.Strength)

			// Update match with results
			match.HomeTeamScore = homeGoals
			match.AwayTeamScore = awayGoals
			match.IsPlayed = true

			// Update match in database
			if err := s.matchRepo.Update(&match); err != nil {
				return nil, err
			}

			// Update home team stats
			homeTeam.Stats.GoalsFor += homeGoals
			homeTeam.Stats.GoalsAgainst += awayGoals

			// Update away team stats
			awayTeam.Stats.GoalsFor += awayGoals
			awayTeam.Stats.GoalsAgainst += homeGoals

			// Update points and win/draw/loss records
			if homeGoals > awayGoals {
				// Home team wins
				homeTeam.Stats.Points += 3
				homeTeam.Stats.Wins++
				awayTeam.Stats.Losses++
			} else if homeGoals < awayGoals {
				// Away team wins
				awayTeam.Stats.Points += 3
				awayTeam.Stats.Wins++
				homeTeam.Stats.Losses++
			} else {
				// Draw
				homeTeam.Stats.Points++
				awayTeam.Stats.Points++
				homeTeam.Stats.Draws++
				awayTeam.Stats.Draws++
			}

			// Save updated teams
			if err := s.teamRepo.Update(homeTeam); err != nil {
				return nil, err
			}
			if err := s.teamRepo.Update(awayTeam); err != nil {
				return nil, err
			}
		}
	}

	// Get all teams for the league table
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Sort teams by points (descending) and goal difference (descending)
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Stats.Points != teams[j].Stats.Points {
			return teams[i].Stats.Points > teams[j].Stats.Points
		}
		// Calculate goal difference on the fly
		iGoalDiff := teams[i].Stats.GoalsFor - teams[i].Stats.GoalsAgainst
		jGoalDiff := teams[j].Stats.GoalsFor - teams[j].Stats.GoalsAgainst
		return iGoalDiff > jGoalDiff
	})

	return teams, nil
}
