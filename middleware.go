package main

import "net/http"

func (cfg *apiConfig) middlewearMetricsInc(next http.Handler) http.HandlerFunc {
	newHandler := func(r http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(r, req)
	}

	return newHandler
}
