package main

import (
	"net/http"

	"github.com/EthanJacksonProject/chirpy/internal/database"
)

func DatabaseChirpToChirp(dbChirp database.Chirp) Chirp {
	newChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	return newChirp
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// 	[
	//   {
	//     "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
	//     "created_at": "2021-01-01T00:00:00Z",
	//     "updated_at": "2021-01-01T00:00:00Z",
	//     "body": "Yo fam this feast is lit ong",
	//     "user_id": "123e4567-e89b-12d3-a456-426614174000"
	//   },

	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to get chirps from DB")
		return
	}

	chirpsSlice := make([]Chirp, len(chirps))

	for i, v := range chirps {
		chirpsSlice[i] = DatabaseChirpToChirp(v)
	}

	// resp, err := json.Marshal(chirpsSlice)
	// if err != nil {
	// 	respondWithError(w, 500, "Failed to marshal to JSON")
	// 	return
	// }

	respondWithJSON(w, 200, chirpsSlice)

}
