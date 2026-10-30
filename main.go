package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func readynessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) adminMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	content := fmt.Sprintf(
		`<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(content))

}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) validateHandler(w http.ResponseWriter, r *http.Request) {
	//"body": "This is an opinion I need to share with the world"'
	type parameters struct {
		Body string `json:"body"`
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

	// it passses
	// create a Valid struct to represent basic success message
	type Valid struct {
		Res bool `json:"valid"`
	}
	resp := Valid{true}
	respondWithJSON(w, 200, resp)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnVals struct {
		ErrMsg string `json:"error"`
	}

	respBody := returnVals{ErrMsg: msg}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(dat))
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	// Take any interface and convert it to JSON
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func (cfg *apiConfig) middlewearMetricsInc(next http.Handler) http.HandlerFunc {
	newHandler := func(r http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(r, req)
	}

	return newHandler
}

func main() {
	mux := http.NewServeMux()

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	apiCfg := apiConfig{}
	appMiddlewareHandler := apiCfg.middlewearMetricsInc(appHandler)

	mux.Handle("/app/", appMiddlewareHandler)
	mux.HandleFunc("GET /api/healthz", readynessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminMetricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateHandler)

	server := http.Server{Handler: mux, Addr: ":8080"}

	server.ListenAndServe()
}
