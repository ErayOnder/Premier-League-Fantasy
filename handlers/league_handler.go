package handlers

import (
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
	teams, err := h.service.PlayWeek(week)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(teams)
}
