package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByID(writer http.ResponseWriter, request *http.Request) {
	userID, err := checkAuth(request.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(writer, 401, "Invalid auth token", err)
		return
	}
	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, 400, "Invalid or missing chirpID", err)
		return
	}

	log.Printf("Chirp ID is %v", chirpID)
	dbChirp, err := cfg.db.GetChirpByID(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 404, "Chirp not found in the database", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(writer, 403, "Forbidden", err)
		return
	}

	err = cfg.db.DeleteChirpByID(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 500, "Failed to delete Chirp", err)
		return
	}

	respondWithJSON(writer, 204, nil)
}
