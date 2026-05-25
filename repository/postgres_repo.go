package repository

import (
	"database/sql"
	"insider-league/domain"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func (r *PostgresRepo) GetTeams() ([]domain.Team, error) {
	rows, err := r.DB.Query("SELECT id, name, strength FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		rows.Scan(&t.ID, &t.Name, &t.Strength)
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *PostgresRepo) GetMatches() ([]domain.Match, error) {
	rows, err := r.DB.Query("SELECT id, home_team_id, away_team_id, week, home_score, away_score, is_played FROM matches ORDER BY week ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		rows.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.Week, &m.HomeScore, &m.AwayScore, &m.IsPlayed)
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *PostgresRepo) UpdateMatch(matchID, homeScore, awayScore int) error {
	_, err := r.DB.Exec("UPDATE matches SET home_score = $1, away_score = $2, is_played = TRUE WHERE id = $3", homeScore, awayScore, matchID)
	return err
}

func (r *PostgresRepo) GetMatchesByWeek(week int) ([]domain.Match, error) {
	rows, err := r.DB.Query("SELECT id, home_team_id, away_team_id, week, home_score, away_score, is_played FROM matches WHERE week = $1 AND is_played = true ORDER BY id ASC", week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		err := rows.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.Week, &m.HomeScore, &m.AwayScore, &m.IsPlayed)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}
