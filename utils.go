package main

import (
	"net/http"

	"Chirpy/internal/auth"

	"github.com/google/uuid"
)

func checkAuth(headers http.Header, secret string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(headers)
	if err != nil || token == "" {
		return uuid.New(), err
	}
	userID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		return uuid.New(), err
	}
	return userID, nil
}
