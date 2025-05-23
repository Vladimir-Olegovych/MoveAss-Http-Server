package handlers

import (
	"encoding/json"
	"goproject/internal/storage"
	"goproject/internal/token"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Processor struct {
	Pages           [2]*template.Template
	TokenService    token.TokenService
	DataBaseService storage.DataBaseService
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserStatsModel struct {
	Money int64 `json:"money"`
}

type UserDataModel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p *Processor) MainHandler(w http.ResponseWriter, r *http.Request) {
	if err := p.Pages[0].Execute(w, nil); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func (p *Processor) SecondHandler(w http.ResponseWriter, r *http.Request) {
	if err := p.Pages[1].Execute(w, nil); err != nil {
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

	model, err := p.TokenService.ParseToken(token)
	if err != nil {
		log.Printf("token error: %v", err)
		http.Error(w, "Expired token not found", http.StatusNotAcceptable)
		return
	}

	userStats, err := p.DataBaseService.FindStatById(model.UUID)
	if err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, "Stats not found", http.StatusNotFound)
		return
	}

	resp := UserStatsModel{
		Money: userStats.UserMoney,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func (p *Processor) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var data UserDataModel
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !validateUserData(data) {
		http.Error(w, "Missing name or password", http.StatusBadRequest)
		return
	}

	userModel, err := p.DataBaseService.FindUser(data.Name, data.Password)
	if err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tokenStr, err := p.TokenService.GenerateToken(userModel.ID)
	if err != nil {
		log.Printf("token generation error: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token: tokenStr,
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

	if value, _ := p.DataBaseService.FindUser(data.Name, data.Password); value != nil {
		http.Error(w, "User has already created", http.StatusConflict)
		return
	}

	id := uuid.NewString()

	tokenStr, err := p.TokenService.GenerateToken(id)
	if err != nil {
		log.Printf("token generation error: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	p.DataBaseService.CreateUser(id, data.Name, data.Password)

	resp := LoginResponse{
		Token: tokenStr,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
func validateUserData(data UserDataModel) bool {
	return data.Name != "" && data.Password != ""
}
