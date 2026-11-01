package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/EthanJacksonProject/chirpy/internal/database"
	"github.com/google/uuid"
)

// func censorWords(content string) string {
// 	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
// 	res := ""

// 	for _, wrd := range strings.Split(content, " ") {
// 		if slices.Contains(bad_words, strings.ToLower(wrd)) {
// 			res = res + " ****"
// 		} else {
// 			res = res + " " + wrd
// 		}
// 	}

// 	return res[1:]
// }

// provided solution that skips weird string trimming
func censorWords(content string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(content, " ")
	for i, wrd := range words {
		if _, ok := badWords[strings.ToLower(wrd)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {
	//"body": "This is an opinion I need to share with the world"'
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// something went wrong when decoding POST
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// Validate chirp length
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// Censor cursing
	cleaned_msg := censorWords(params.Body)

	// it passses
	// create a Valid struct to represent basic success message
	type Valid struct {
		Res string `json:"cleaned_body"`
	}
	validRes := Valid{cleaned_msg}

	// .Body, arg.UserID
	type chirpParams struct {
		Body   string
		UserID uuid.UUID
	}

	chirpy := database.CreateChirpParams{
		Body:   validRes.Res,
		UserID: params.UserID,
	}
	chirpRes, err := cfg.dbQueries.CreateChirp(r.Context(), chirpy)

	friendlyChirp := Chirp{
		ID:        chirpRes.ID,
		CreatedAt: chirpRes.CreatedAt,
		UpdatedAt: chirpRes.UpdatedAt,
		Body:      chirpRes.Body,
		UserID:    chirpRes.UserID,
	}
	respondWithJSON(w, 201, friendlyChirp)
}
