package main

import (
	"log"
	"net/http"
	"time"

	"Chirpy/internal/auth"
	"Chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRefreshToken(writer http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 401, "Refresh token invalid or not provided", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, 401, "Refresh token not found or expired", err)
		return
	}

	expirationTime := time.Hour
	token, err := auth.MakeJWT(userID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(writer, 500, "Failed in creating jwt token", err)
		return
	}

	type TokenBody struct {
		Token string `json:"token"`
	}

	tokenBody := TokenBody{
		Token: token,
	}

	respondWithJSON(writer, 200, tokenBody)
}

func (cfg *apiConfig) handlerRevokeRefreshToken(writer http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 500, "Refresh token invalid or not provided", err)
		return
	}

	err = cfg.db.RevokeRefreshTokenByToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, 500, "Refresh token not found in the database", err)
		return
	}

	respondWithJSON(writer, 204, nil)
}

func (cfg *apiConfig) CreateAndStoreRefreshToken(userID uuid.UUID, request *http.Request) string {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Failed to create refrest token with the error: %v", err)
		return ""
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: userID,
	}
	returnedToken, err := cfg.db.CreateRefreshToken(request.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Failed to store refresh token in the database with error: %v", err)
		return ""
	}

	return returnedToken.Token
}
