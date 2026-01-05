package main

import (
	"log"
	"net/http"
	"sort"

	"Chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetAll(writer http.ResponseWriter, request *http.Request) {
	authorID := request.URL.Query()["author_id"]
	var err error
	dbChirps := make([]database.Chirp, 0, 100)
	if authorID == nil {

		dbChirps, err = cfg.db.GetAllChirps(request.Context())
		if err != nil {
			respondWithError(writer, 500, "Unable to get chirps from the db", err)
			return
		}
	} else {

		userID, err := uuid.Parse(authorID[0])
		if err != nil {
			respondWithError(writer, 500, "Unable to get chirps from the db", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUser(request.Context(), userID)
		if err != nil {
			respondWithError(writer, 404, "User not found", err)
			return
		}
	}
	paramSort := request.URL.Query()["sort"]
	log.Printf("query params are: %v", request.URL.Query())
	if paramSort != nil && paramSort[0] == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt) })
	} else {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt) })
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
		respondWithError(writer, 404, "Chirp not found", err)
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
