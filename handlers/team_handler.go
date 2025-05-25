package handlers

import (
	"insider-league/models"
	"insider-league/services"

	"github.com/gofiber/fiber/v2"
)

// TeamHandler handles team-related HTTP requests
type TeamHandler struct {
	service services.TeamService
}

// NewTeamHandler creates and returns a new TeamHandler instance
func NewTeamHandler(service services.TeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

// CreateTeam handles the creation of a new team
func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	team := new(models.Team)

	// Parse the request body into the team struct
	if err := c.BodyParser(team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Create the team using the service
	if err := h.service.Create(team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the created team with a 201 status code
	return c.Status(fiber.StatusCreated).JSON(team)
}

// GetAllTeams handles retrieving all teams
func (h *TeamHandler) GetAllTeams(c *fiber.Ctx) error {
	teams, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(teams)
}
