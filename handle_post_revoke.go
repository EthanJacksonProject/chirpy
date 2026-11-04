package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EthanJacksonProject/chirpy/internal/auth"
)

// type User struct {
// 	ID           uuid.UUID `json:"id"`
// 	CreatedAt    time.Time `json:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at"`
// 	Email        string    `json:"email"`
// 	Token        string    `json:"token"`
// 	RefreshToken string    `json:"refresh_token"`
// }

// func UserFromDBUser(dbUsr database.User) User {
// 	return User{
// 		ID:        dbUsr.ID,
// 		Email:     dbUsr.Email,
// 		CreatedAt: dbUsr.CreatedAt,
// 		UpdatedAt: dbUsr.UpdatedAt,
// 	}
// }

func (cfg *apiConfig) postRevokeHandler(w http.ResponseWriter, r *http.Request) {
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

	if dbRefreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Expired refresh: %v", err))
		return
	}
	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Token revoked: %v", err))
		return
	}

	if dbRefreshToken.Token == "" {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Refresh Token Empty: %v", err))
		return
	}

	//
	err = cfg.dbQueries.RevokeTokenByToken(r.Context(), dbRefreshToken.Token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Failed to revoke token: %v", err))
		return
	}

	// type Empty struct {
	// }
	// respondWithJSON(w, http.StatusNoContent, Empty{})
	w.WriteHeader(http.StatusNoContent)
	return

}
