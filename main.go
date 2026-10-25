package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func readynessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) apiHitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))

}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

type apiConfig struct {
	fileserverHits atomic.Int32
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
	mux.HandleFunc("GET /healthz", readynessHandler)
	mux.HandleFunc("GET /metrics", apiCfg.apiHitHandler)
	mux.HandleFunc("POST /reset", apiCfg.resetHandler)

	server := http.Server{Handler: mux, Addr: ":8080"}

	server.ListenAndServe()
}
