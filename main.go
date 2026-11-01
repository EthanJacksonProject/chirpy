package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/EthanJacksonProject/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	godotenv.Load()

	// Set Up Database Stuff
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("DB Connection Failed")
		return
	}

	// serve route any paths starting with `/app/` to serve out of this directory
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	// set up shared API config
	apiCfg := apiConfig{}

	//dbQueries:= database.New(db)
	// Hook up DB querier to api config
	apiCfg.dbQueries = database.New(db)

	// set platform
	apiCfg.platform = os.Getenv("PLATFORM")

	// Hook up middleware hit counter
	// any hits to `/app/` will increment
	appMiddlewareHandler := apiCfg.middlewearMetricsInc(appHandler)

	// create root mux for routing
	mux := http.NewServeMux()
	mux.Handle("/app/", appMiddlewareHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.adminMetricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("GET /api/healthz", readynessHandler)
	mux.HandleFunc("POST /api/users", apiCfg.postUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirpHandler)

	server := http.Server{Handler: mux, Addr: ":8080"}

	server.ListenAndServe()
}
