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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	expirationTime := time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(writer, 500, "Failed in creating jwt token", err)
	}

	refreshToken := cfg.CreateAndStoreRefreshToken(user.ID, request)
	userObj := DBUserToUser(user, token, refreshToken)
	userObj.Token = token

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

	user, err := cfg.db.GetUserByEmail(request.Context(), reqBody.Email)
	if err != nil {
		respondWithError(writer, 500, "User not found", err)
		return
	}

	success, err := auth.CheckPasswordAndHash(reqBody.Password, user.HashedPassword)
	if err != nil || !success {
		respondWithError(writer, 401, "Password not valid", err)
		return
	}

	expirationTime := time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(writer, 500, "Could not create JWT Token", err)
		return
	}

	refreshToken := cfg.CreateAndStoreRefreshToken(user.ID, request)
	userObj := DBUserToUser(user, token, refreshToken)
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

func (cfg *apiConfig) handlerUserModification(writer http.ResponseWriter, request *http.Request) {
	userID, err := checkAuth(request.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(writer, 401, "Invalid or missing auth token", err)
	}

	reqBody := reqUserBody{}
	err = json.NewDecoder(request.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(writer, 500, "Failed to decode request body with error", err)
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		respondWithError(writer, 400, "User change request requires an email and a password:", err)
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(writer, 500, "Failed to hash password", err)
	}

	updateUserInfoParams := database.UpdateUserInfoParams{
		ID:             userID,
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.UpdateUserInfo(request.Context(), updateUserInfoParams)
	if err != nil {
		respondWithError(writer, 500, "Error updating the database", err)
	}

	returnBody := DBUserToUser(user, "", "")
	respondWithJSON(writer, 200, returnBody)
}
