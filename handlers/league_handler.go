package handlers

import (
	"fmt"
	"insider-league/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// LeagueHandler handles league-related HTTP requests
type LeagueHandler struct {
	service services.LeagueService
}

// NewLeagueHandler creates and returns a new LeagueHandler instance
func NewLeagueHandler(service services.LeagueService) *LeagueHandler {
	return &LeagueHandler{
		service: service,
	}
}

// PlayWeek handles simulating all matches for a given week
func (h *LeagueHandler) PlayWeek(c *fiber.Ctx) error {
	// Get and parse the week parameter
	weekStr := c.Params("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}

	// Play the week's matches
	teams, matches, predictions, err := h.service.PlayWeek(week)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// If teams is nil, it means the week was already played
	if teams == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": fmt.Sprintf("Week %d has already been played", week),
			"matches": matches,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"league_table": teams,
		"matches":      matches,
		"predictions":  predictions,
	})
}

// PlayAll handles simulating all remaining matches in the league
func (h *LeagueHandler) PlayAll(c *fiber.Ctx) error {
	// Play all remaining matches
	leagueTable, err := h.service.PlayAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"league_table": leagueTable,
	})
}

// EditMatchResult handles updating a match result
func (h *LeagueHandler) EditMatchResult(c *fiber.Ctx) error {
	// Get match ID from URL parameters
	matchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	// Parse request body
	type updateScoreRequest struct {
		HomeGoals int `json:"home_goals"`
		AwayGoals int `json:"away_goals"`
	}

	var req updateScoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate scores
	if req.HomeGoals < 0 || req.AwayGoals < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Scores cannot be negative",
		})
	}

	// Update match result
	match, leagueTable, err := h.service.EditMatchResult(matchID, req.HomeGoals, req.AwayGoals)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"match":        match,
		"league_table": leagueTable,
	})
}

// ResetLeague resets all match results and team statistics
func (h *LeagueHandler) ResetLeague(c *fiber.Ctx) error {
	if err := h.service.ResetLeague(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "League has been reset successfully",
	})
}

// GetLeagueTable retrieves the current league table
func (h *LeagueHandler) GetLeagueTable(c *fiber.Ctx) error {
	teams, err := h.service.GetLeagueTable()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"teams": teams,
	})
}
