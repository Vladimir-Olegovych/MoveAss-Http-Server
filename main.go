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
	pages := [2]*template.Template{
		template.Must(template.ParseFiles("res/templates/index.html")),
		template.Must(template.ParseFiles("res/templates/index_second.html")),
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
		Pages:           pages,
		TokenService:    tokenService,
		DataBaseService: databaseService,
	}

	defer databaseService.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", p.MainHandler)
	mux.HandleFunc("GET /second", p.SecondHandler)
	mux.HandleFunc("POST /create", p.CreateUserHandler)
	mux.HandleFunc("POST /login", p.LoginUserHandler)
	mux.HandleFunc("GET /stats", p.GetUserStatsHandler)

	port := ":5123"
	log.Printf("Starting server on %s", port)
	serverError := http.ListenAndServe(port, mux)
	return serverError
}
