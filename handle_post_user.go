package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
	"github.com/EthanJacksonProject/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func UserFromDBUser(dbUsr database.User) User {
	return User{
		ID:          dbUsr.ID,
		Email:       dbUsr.Email,
		CreatedAt:   dbUsr.CreatedAt,
		UpdatedAt:   dbUsr.UpdatedAt,
		IsChirpyRed: dbUsr.IsChirpyRed,
	}
}

func (cfg *apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
	// unmarshal body to get email and responds with User struct
	type UserBodyRaw struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userBody := UserBodyRaw{}
	err := decoder.Decode(&userBody)
	// something went wrong when decoding POST
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	hashedPW, _ := auth.HashPassword(userBody.Password)

	userParams := database.CreateUserParams{
		Email:          userBody.Email,
		HashedPassword: hashedPW,
	}

	dbUsr, err := cfg.dbQueries.CreateUser(r.Context(), userParams)

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
