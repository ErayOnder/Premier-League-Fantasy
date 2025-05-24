package main

import (
	"fmt"
	"insider-league/db"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Check required environment variables
	requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "DB_PORT", "DB_SSLMODE"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	// Connect to the database
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Premier League Fantasy API")
	})

	// API routes
	api := app.Group("/api")

	// Teams routes
	teams := api.Group("/teams")
	teams.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Will get all teams")
	}) // For now, just return a string

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
