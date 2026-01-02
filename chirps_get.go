package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetAll(writer http.ResponseWriter, request *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(request.Context())
	if err != nil {
		respondWithError(writer, 500, "Unable to get chirps from the db", err)
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		}
		chirps = append(chirps, chirp)
	}
	respondWithJSON(writer, 200, chirps)
}

func (cfg *apiConfig) handlerChirpsGetByID(writer http.ResponseWriter, request *http.Request) {
	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, 400, "Invalid UUID", err)
	}

	dbChirp, err := cfg.db.GetChirpByID(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 404, "Chip not found", err)
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
	}
	respondWithJSON(writer, 200, chirp)
}
