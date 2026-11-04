package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpsByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), parsedID)

	if err != nil {
		respondWithError(w, 404, "Err: Chirp not found")
		return
	}

	chirp := DatabaseChirpToChirp(dbChirp)
	respondWithJSON(w, 200, chirp)
}
