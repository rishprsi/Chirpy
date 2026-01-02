package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Chirpy/internal/auth"
	"Chirpy/internal/database"

	"github.com/google/uuid"
)

type reqUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) CreateUserHandler(writer http.ResponseWriter, request *http.Request) {
	reqBody := reqUserBody{}
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(writer, 500, "Could not decode request body", err)
		return
	}
	if reqBody.Email == "" || reqBody.Password == "" {
		respondWithError(writer, 400, "No uslr email provided, please provide user email\n", nil)
		return
	}
	hash, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(writer, 500, "Failed to compute hash", err)
	}
	params := database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hash,
	}
	user, err := cfg.db.CreateUser(request.Context(), params)
	if err != nil {
		respondWithError(writer, 500, "Error creating user", err)
		return
	}

	userObj := DBUserToUser(user)

	log.Printf("Created the user with the email %v", user.Email)
	respondWithJSON(writer, 201, userObj)
}

func (cfg *apiConfig) handlerUserLogin(writer http.ResponseWriter, request *http.Request) {
	reqBody := reqUserBody{}
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(writer, 500, "Failed to decode request body", err)
		return
	}
	if reqBody.Email == "" || reqBody.Password == "" {
		respondWithError(writer, 500, "Email and Password are required for login", nil)
		return
	}

	user, err := cfg.db.GetUser(request.Context(), reqBody.Email)
	if err != nil {
		respondWithError(writer, 500, "User not found", err)
		return
	}

	success, err := auth.CheckPasswordAndHash(reqBody.Password, user.HashedPassword)
	if err != nil || !success {
		respondWithError(writer, 401, "Password not valid", err)
		return
	}

	userObj := DBUserToUser(user)
	respondWithJSON(writer, 200, userObj)
}

func (cfg *apiConfig) ResetHandler(writer http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(writer, 403, "Forbidden outside dev", nil)
	}
	err := cfg.db.DeleteUsers(request.Context())
	if err != nil {
		respondWithError(writer, 500, "Failed to delete users", err)
	}
}
