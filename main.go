package main

import (
	"fmt"
	"insider-league/db"
	"insider-league/db/seeds"
	"insider-league/handlers"
	"insider-league/repository"
	"insider-league/services"
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

	// Seed the database
	if err := seeds.Load(db.DB); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize repositories
	teamRepo := repository.NewTeamRepository(db.DB)
	matchRepo := repository.NewMatchRepository(db.DB)

	// Initialize services
	teamService := services.NewTeamService(teamRepo)
	matchService := services.NewMatchService(matchRepo)
	leagueService := services.NewLeagueService(teamService, matchService)

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
	teamHandler := handlers.NewTeamHandler(teamService)
	teams.Get("/", teamHandler.GetAllTeams)
	teams.Get("/:id", teamHandler.GetTeamByID)
	teams.Put("/:id", teamHandler.UpdateTeam)
	teams.Delete("/:id", teamHandler.DeleteTeam)
	teams.Post("/", teamHandler.CreateTeam)

	// Matches routes
	matches := api.Group("/matches")
	matchHandler := handlers.NewMatchHandler(matchService)
	matches.Get("/", matchHandler.GetAllMatches)
	matches.Get("/:id", matchHandler.GetMatchByID)
	matches.Put("/:id", matchHandler.UpdateMatch)
	matches.Delete("/:id", matchHandler.DeleteMatch)
	matches.Post("/", matchHandler.CreateMatch)

	// League routes
	league := api.Group("/league")
	leagueHandler := handlers.NewLeagueHandler(leagueService)
	league.Get("/", leagueHandler.GetLeagueTable)
	league.Get("/play", leagueHandler.PlayNextWeek)
	league.Get("/play-all", leagueHandler.PlayAll)
	league.Get("/week/:id", leagueHandler.GetWeekResults)
	league.Put("/edit-match/:id", leagueHandler.EditMatchResult)
	league.Post("/reset", leagueHandler.ResetLeague)

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
