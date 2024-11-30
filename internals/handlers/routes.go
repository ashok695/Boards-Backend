package handlers

import (
	"boards/internals/services" // Import services
	"github.com/gofiber/fiber/v2"
)

// GetBoardDetailsHandler handles requests to fetch board details
func GetBoardDetailsHandler(c *fiber.Ctx) error {
	// Get board details from the service
	boardData, err := services.GetBoardDetails()
	if err != nil {
		// If there's an error, return a 500 response
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the board data as a JSON response
	return c.Status(200).JSON(boardData)
}
func GetBoardTasksHandler(c *fiber.Ctx) error {
	boardData, err := services.GetBoardTasksDetails()
	if err != nil {
		// If there's an error, return a 500 response
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the board data as a JSON response
	return c.Status(200).JSON(boardData)
}