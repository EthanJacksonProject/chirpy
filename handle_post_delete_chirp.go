package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to fetch token")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to validate token")
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirp")
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you cannot delete this chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
