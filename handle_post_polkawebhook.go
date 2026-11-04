package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
	"github.com/EthanJacksonProject/chirpy/internal/database"
	"github.com/google/uuid"
)

type PolkaData struct {
	UserID string `json:"user_id"`
}

type PolkaParams struct {
	Event string    `json:"event"`
	Data  PolkaData `json:"data"`
}

func (cfg *apiConfig) postPolkaWebhookHandler(w http.ResponseWriter, r *http.Request) {

	authToken, err := auth.GetAPIKey(r.Header)
	if authToken != cfg.polkaKey || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid auth token")
		return
	}

	params := PolkaParams{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbParams := database.SetChirpyRedStatusByUserIDParams{
		ID:          userID,
		IsChirpyRed: true,
	}
	_, err = cfg.dbQueries.SetChirpyRedStatusByUserID(r.Context(), dbParams)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
