package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
)

func (cfg *apiConfig) postRefreshHandler(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("failed to get bearer token from header: %v", err))
		return
	}

	dbRefreshToken, err := cfg.dbQueries.GetTokenDataByToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("failed to retrieve refresh token from db: %v", err))
		return
	}

	if dbRefreshToken.ExpiresAt.Before(time.Now().UTC()) || dbRefreshToken.RevokedAt.Valid || dbRefreshToken.Token == "" {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid refresh token: %v", err))
		return
	}

	type Resp struct {
		Token string `json:"token"`
	}

	newToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.jwtSecret, time.Duration(60*time.Minute))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to refresh JWT: %v", err))
		return
	}
	res := Resp{Token: newToken}
	respondWithJSON(w, http.StatusOK, res)

}
