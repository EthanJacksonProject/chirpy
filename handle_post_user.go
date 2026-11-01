package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
	// unmarshal body to get email and responds with User struct
	type UserEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userEmail := UserEmail{}
	err := decoder.Decode(&userEmail)

	// something went wrong when decoding POST
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	email := userEmail.Email

	dbUsr, err := cfg.dbQueries.CreateUser(r.Context(), email)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// DB User has the wrong labels, as they match the SQL column names
	// So, we initialize a User{} with the same data so match what our API wants
	// This copies the sql column names over into our friendly Go struct
	usr := User{ID: dbUsr.ID,
		Email:     dbUsr.Email,
		CreatedAt: dbUsr.CreatedAt,
		UpdatedAt: dbUsr.UpdatedAt}

	respondWithJSON(w, 201, usr)

}
