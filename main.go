package main

import (
	"goproject/internal/handlers"
	"goproject/internal/storage"
	"goproject/internal/token"
	"html/template"
	"log"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func run() error {
	html, templateError := template.ParseFiles("res/templates/index.html")
	if templateError != nil {
		return templateError
	}
	databaseService, databaseError := storage.Open(
		"res/sql/store.db",
		"res/sql/users.sql",
		"res/sql/stats.sql",
	)
	if databaseError != nil {
		return databaseError
	}
	tokenService := &token.JWTService{}

	p := handlers.Processor{
		Page:            html,
		TokenService:    tokenService,
		DataBaseService: databaseService,
	}

	defer databaseService.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", p.IndexHandler)
	mux.HandleFunc("POST /create", p.CreateUserHandler)
	mux.HandleFunc("POST /stats", p.GetUserStatsHandler)

	port := ":8080"
	log.Printf("Starting server on %s", port)
	serverError := http.ListenAndServe(port, mux)
	return serverError
}
