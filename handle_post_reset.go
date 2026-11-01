package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if cfg.platform != "dev" {
		respondWithError(w, 403, "Platform != dev; reset forbidden")
		return
	}

	cfg.dbQueries.DropAllUsers(r.Context())
}
