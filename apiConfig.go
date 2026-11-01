package main

import (
	"sync/atomic"

	"github.com/EthanJacksonProject/chirpy/internal/database"
)

// Shared values for all API calls
type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}
