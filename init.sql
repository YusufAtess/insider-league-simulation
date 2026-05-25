CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    strength INT NOT NULL
);

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team_id INT REFERENCES teams(id),
    away_team_id INT REFERENCES teams(id),
    week INT NOT NULL,
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    is_played BOOLEAN DEFAULT FALSE
);

INSERT INTO teams (id, name, strength) VALUES 
(1, 'Chelsea', 85), (2, 'Arsenal', 75), (3, 'Manchester City', 80), (4, 'Liverpool', 70);

INSERT INTO matches (home_team_id, away_team_id, week) VALUES (1, 2, 1), (3, 4, 1);
INSERT INTO matches (home_team_id, away_team_id, week) VALUES (4, 1, 2), (2, 3, 2);
INSERT INTO matches (home_team_id, away_team_id, week) VALUES (1, 3, 3), (4, 2, 3);
INSERT INTO matches (home_team_id, away_team_id, week) VALUES (2, 1, 4), (4, 3, 4);
INSERT INTO matches (home_team_id, away_team_id, week) VALUES (1, 4, 5), (3, 2, 5);
INSERT INTO matches (home_team_id, away_team_id, week) VALUES (3, 1, 6), (2, 4, 6);