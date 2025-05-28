-- Teams table
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    strength INTEGER NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    goals_for INTEGER NOT NULL DEFAULT 0,
    goals_against INTEGER NOT NULL DEFAULT 0,
    goal_difference INTEGER NOT NULL DEFAULT 0,
    wins INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0
);

-- Matches table
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    week INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_team_score INTEGER NOT NULL DEFAULT 0,
    away_team_score INTEGER NOT NULL DEFAULT 0,
    is_played BOOLEAN NOT NULL DEFAULT false
);

-- Add indexes for better query performance
CREATE INDEX idx_matches_week ON matches(week);
CREATE INDEX idx_matches_home_team_id ON matches(home_team_id);
CREATE INDEX idx_matches_away_team_id ON matches(away_team_id); 