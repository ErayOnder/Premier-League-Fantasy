# Insider Case Study: Premier League Simulation

## Overview

This project is a backend simulation of a 4-team Premier League season, built with Go and Fiber web framework. It tracks match results, updates a league table with Premier League rules, and provides championship predictions. The project fulfills the requirements of a complete season simulation, match management, and real-time league table editions.

## Features

- **Full league simulation** for 4 teams with 6 weeks of matches
- **Premier League rules** implementation for points calculation and table sorting
- **Championship predictions** after week 4 based on current standings
- **Complete API** to manage teams and matches with full CRUD operations
- **"Play All" functionality** to simulate the entire season at once
- **"Play Next Week" functionality** to simulate matches week by week
- **"Edit Match Result" functionality** with automatic league table recalculation
- **Automatic database seeding** with teams and full season fixtures
- **Real-time league standings** with points, goals, and goal difference tracking
- **Week-specific results** viewing for match history

## Tech Stack

- **Go** (1.24.3)
- **Fiber** (Web Framework)
- **GORM** (ORM for database operations)
- **PostgreSQL** (Database)
- **godotenv** (Environment variable management)

## Prerequisites

Before running this project, ensure you have the following installed:

- **Go** (version 1.24.3 or higher)
- **PostgreSQL** (version 12 or higher)
- **Postman** (for API testing)

## Setup & Running the Project

### 1. Clone the repository
```bash
git clone https://github.com/ErayOnder/Premier-League-Fantasy.git
cd insider-league
```

### 2. Set up Environment Variables

Create a `.env` file in the project root with your database configuration. Use the following template:

```env
DB_HOST=localhost
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=insider_league_db
DB_PORT=5432
DB_SSLMODE=disable
SERVER_PORT=8080
```

**Note:** All environment variables listed above are required for the application to start.

### 3. Set up PostgreSQL Database

Ensure PostgreSQL is running and create a database with the name specified in your `.env` file:

```sql
CREATE DATABASE insider_league_db;
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Application

```bash
go run main.go
```

The application will:
- Connect to the database
- Automatically create the required tables
- Seed the database with 4 teams and all season fixtures
- Start the server on the specified port (default: 8080)

You should see output like:
```
Server starting on port 8080
```

## Using the API with Postman

### Import the Postman Collection

1. Open Postman
2. Import the `LeagueEndpoints.postman_collection.json` file located in the project root
3. The collection is organized into folders for easy navigation:
   - **League**: Core simulation endpoints (Play Next Week, Play All, Get League Table, etc.)
   - **Teams**: Team management endpoints
   - **Matches**: Match management endpoints

### Key API Endpoints

The screenshots of the results can be find under `/endpointscreenshots` folder.
#### League Simulation
- `GET /api/league/` - Get current league table/standings
- `GET /api/league/play` - Play the next week's matches
- `GET /api/league/play-all` - Simulate all remaining matches
- `GET /api/league/week/:id` - Get results for a specific week
- `PUT /api/league/edit-match/:id` - Edit a match result (recalculates league table)
- `POST /api/league/reset` - Reset the entire league (clears all match results)

#### Teams
- `GET /api/teams/` - Get all teams
- `GET /api/teams/:id` - Get specific team details
- `POST /api/teams/` - Create a new team
- `PUT /api/teams/:id` - Update team information
- `DELETE /api/teams/:id` - Delete a team

#### Matches
- `GET /api/matches/` - Get all matches
- `GET /api/matches/:id` - Get specific match details
- `POST /api/matches/` - Create a new match
- `PUT /api/matches/:id` - Update match details
- `DELETE /api/matches/:id` - Delete a match

### Typical Usage Flow

1. **View Initial State**: Use `GET /api/league/` to see the initial league table
2. **Simulate Matches**: Use `GET /api/league/play` to play week by week, or `GET /api/league/play-all` to simulate the entire season
3. **Check Results**: Use `GET /api/league/week/:id` to see specific week results
4. **Edit if Needed**: Use `PUT /api/league/edit-match/:id` to modify match results
5. **Reset**: Use `POST /api/league/reset` to start over

## Database Schema

The database schema consists of two main tables and can be found in the `schema.sql` file:

### Teams Table
- Stores team information including name, strength, and league statistics
- Tracks points, goals for/against, goal difference, wins, draws, and losses

### Matches Table
- Stores fixture information and results
- Links to home and away teams
- Tracks week number, scores, and whether the match has been played

## Project Structure

```
insider-league/
├── db/                     # Database connection and seeding
├── handlers/               # HTTP request handlers
├── models/                 # Data models/structs
├── repository/             # Data access layer
├── services/               # Business logic layer
├── helpers/                # Utility functions
├── mocks/                  # Test mocks
├── main.go                 # Application entry point
├── schema.sql              # Database schema
├── go.mod                  # Go module dependencies
└── LeagueEndpoints.postman_collection.json  # API collection
```