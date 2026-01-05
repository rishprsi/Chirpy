package main

import (
	"encoding/json"
	"net/http"

	"Chirpy/internal/auth"

	"github.com/google/uuid"
)

type PolkaBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUpgradeUser(writer http.ResponseWriter, request *http.Request) {
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(writer, 401, "Invalid or missing API Key", err)
		return
	}

	reqBody := PolkaBody{}
	err = json.NewDecoder(request.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(writer, 400, "Failed to decode the request body", err)
		return
	}

	if reqBody.Event == "user.upgraded" {
		userIDString := reqBody.Data.UserID
		userID, err := uuid.Parse(userIDString)
		if err != nil {
			respondWithError(writer, 500, "Not able to decode user ID from request", err)
			return
		}

		err = cfg.db.UpgradeUserToRed(request.Context(), userID)
		if err != nil {
			respondWithError(writer, 404, "User not found in the database", err)
			return
		}
	}

	respondWithJSON(writer, 204, "User upgraded")
}
