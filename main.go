package main

import (
	"database/sql"
	"fmt"
	"insider-league/handler"
	"insider-league/repository"
	"insider-league/service"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://insider:password@localhost:5432/league_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Dependency Injection
	repo := repository.NewPostgresRepo(db)
	simService := service.NewSimService(repo)
	apiHandler := handler.NewApiHandler(simService)

	http.HandleFunc("/api/standings", apiHandler.GetStandings)
	http.HandleFunc("/api/simulate/next", apiHandler.SimulateNextWeek)
	http.HandleFunc("/api/simulate/all", apiHandler.PlayAll)
	http.HandleFunc("/api/match/edit", apiHandler.EditMatch)
	http.HandleFunc("/api/matches", apiHandler.GetMatchesByWeek)

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
