package handlers

import (
	"insider-league/models"
	"insider-league/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

// GetTeamByID handles retrieving a team by its ID
func (h *TeamHandler) GetTeamByID(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	// Get the team using the service
	team, err := h.service.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Team not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(team)
}

// UpdateTeam handles updating an existing team
func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	// Create a new team instance and parse the request body
	team := new(models.Team)
	if err := c.BodyParser(team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Set the ID from the URL parameter
	team.ID = uint(id)

	// Update the team using the service
	if err := h.service.Update(team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(team)
}

// DeleteTeam handles deleting a team
func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	// Delete the team using the service
	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Team not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
