package handlers

import (
	"encoding/json"
	"goproject/internal/storage"
	"goproject/internal/token"
	"html/template"
	"log"
	"net/http"
)

type Processor struct {
	Page            *template.Template
	TokenService    token.TokenService
	DataBaseService storage.DataBaseService
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserStatsModel struct {
	Name  string `json:"name"`
	Money int64  `json:"money"`
}

type UserDataModel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p *Processor) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := p.Page.Execute(w, nil); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Processor) GetUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token parameter", http.StatusBadRequest)
		return
	}

	userData, err := p.DataBaseService.FindUserByToken(token)
	if err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userStats, err := p.DataBaseService.FindStatsByToken(token)
	if err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, "Stats not found", http.StatusNotFound)
		return
	}

	resp := UserStatsModel{
		Name:  userData.UserName,
		Money: userStats.UserMoney,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func (p *Processor) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var data UserDataModel
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !validateUserData(data) {
		http.Error(w, "Missing name or password", http.StatusBadRequest)
		return
	}

	if value, _ := p.DataBaseService.FindUserByName(data.Name); value != nil {
		http.Error(w, "User has already created", http.StatusConflict)
		return
	}

	tokenStr, err := p.TokenService.GenerateToken(data.Name, "user")
	if err != nil {
		log.Printf("token generation error: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	p.DataBaseService.AddUser(data.Name, data.Password, tokenStr)
	p.DataBaseService.AddStats(tokenStr, 100)

	resp := TokenResponse{Token: tokenStr}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
func validateUserData(data UserDataModel) bool {
	return data.Name != "" && data.Password != ""
}
