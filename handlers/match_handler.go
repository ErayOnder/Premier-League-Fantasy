package handlers

import (
	"insider-league/models"
	"insider-league/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

// GetAllMatches handles retrieving all matches
func (h *MatchHandler) GetAllMatches(c *fiber.Ctx) error {
	matches, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(matches)
}

// GetMatchByID handles retrieving a match by its ID
func (h *MatchHandler) GetMatchByID(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	// Get the match using the service
	match, err := h.service.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Match not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(match)
}
