package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
	"github.com/EthanJacksonProject/chirpy/internal/database"
)

func (cfg *apiConfig) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	thisBody := Body{}
	err := decoder.Decode(&thisBody)
	if err != nil {
		respondWithError(w, 500, "Failed to parse body")
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), thisBody.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	hashMatched, err := auth.CheckPasswordHash(thisBody.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	if !hashMatched {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	usr := UserFromDBUser(dbUser)

	// expire JWT after 1 hour
	token, err := auth.MakeJWT(usr.ID, cfg.jwtSecret, 60*time.Minute)
	if err != nil {
		respondWithError(w, 500, "failed to get bearer token")
		return
	}
	usr.Token = token

	refreshParams := database.CreateRefreshTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: usr.ID,
	}

	refreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	usr.RefreshToken = refreshToken.Token

	respondWithJSON(w, 200, usr)

}
