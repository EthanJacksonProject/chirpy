package main

import (
	"net/http"
	"sort"

	"github.com/EthanJacksonProject/chirpy/internal/database"
	"github.com/google/uuid"
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

	var dbChirps []database.Chirp
	var err error
	var author uuid.UUID

	authorID := r.URL.Query().Get("author_id")

	sortDir := r.URL.Query().Get("sort")
	if sortDir == "" {
		sortDir = "asc"
	}
	if sortDir != "asc" && sortDir != "desc" {
		respondWithError(w, http.StatusBadRequest, "Malformed sort")
		return
	}

	if authorID != "" {
		// Author ID provided
		// look it up
		author, err = uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		dbChirps, err = cfg.dbQueries.GetChirpsByUserID(r.Context(), author)

	} else {
		// Author ID not provided
		// list everything
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, v := range dbChirps {
		chirps[i] = DatabaseChirpToChirp(v)
	}

	// resp, err := json.Marshal(chirpsSlice)
	// if err != nil {
	// 	respondWithError(w, 500, "Failed to marshal to JSON")
	// 	return
	// }

	if sortDir == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			if sortDir == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt) // i was later, after j, so it descends
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) // i was before j, so keep it asending through list
		})
	}

	respondWithJSON(w, 200, chirps)

}
