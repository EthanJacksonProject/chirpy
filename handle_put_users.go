package main

import (
	"encoding/json"
	"net/http"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
	"github.com/EthanJacksonProject/chirpy/internal/database"
)

func (cfg *apiConfig) putUsersHandler(w http.ResponseWriter, r *http.Request) {
	// check auth
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	// user is valid, parse body

	type RequestBody struct {
		Password string
		Email    string
	}

	decoder := json.NewDecoder(r.Body)
	params := RequestBody{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbParams := database.SetEmailAndPWHashByIdParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPW,
	}

	dbUser, err := cfg.dbQueries.SetEmailAndPWHashById(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	usr := UserFromDBUser(dbUser)

	respondWithJSON(w, http.StatusOK, usr)

}
