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

// UpdateMatch handles updating an existing match
func (h *MatchHandler) UpdateMatch(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	// First get the existing match to check if it's played
	existingMatch, err := h.service.GetByID(id)
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

	// Check if the match is already played
	if existingMatch.IsPlayed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot update a match that has already been played",
		})
	}

	// Create a new match instance and parse the request body
	match := new(models.Match)
	if err := c.BodyParser(match); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Set the ID from the URL parameter
	match.ID = uint(id)

	// Update the match using the service
	if err := h.service.Update(match); err != nil {
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

// DeleteMatch handles deleting a match
func (h *MatchHandler) DeleteMatch(c *fiber.Ctx) error {
	// Get and parse the ID parameter
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	// Delete the match using the service
	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Match not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
