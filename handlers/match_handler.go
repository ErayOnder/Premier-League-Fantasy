package handlers

import (
	"insider-league/models"
	"insider-league/services"

	"github.com/gofiber/fiber/v2"
)

// MatchHandler handles match-related HTTP requests
type MatchHandler struct {
	service services.MatchService
}

// NewMatchHandler creates and returns a new MatchHandler instance
func NewMatchHandler(service services.MatchService) *MatchHandler {
	return &MatchHandler{
		service: service,
	}
}

// CreateMatch handles the creation of a new match
func (h *MatchHandler) CreateMatch(c *fiber.Ctx) error {
	match := new(models.Match)

	// Parse the request body into the match struct
	if err := c.BodyParser(match); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Create the match using the service
	if err := h.service.Create(match); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the created match with a 201 status code
	return c.Status(fiber.StatusCreated).JSON(match)
}
