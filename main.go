package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"Chirpy/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file with the error: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to create a connection with the database with err: %v\n", err)
	}

	dbQueries := database.New(db)

	port := "8080"
	filepathRoot := "."
	serveMux := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      os.Getenv("JWT_TOKEN"),
	}

	serveMux.Handle("/app", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	serveMux.HandleFunc("GET /api/healthz", Readiness)
	// serveMux.HandleFunc("POST /api/validate_chirp", ChirpValidation)
	serveMux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	serveMux.HandleFunc("PUT /api/users", cfg.handlerUserModification)
	serveMux.HandleFunc("POST /api/login", cfg.handlerUserLogin)
	serveMux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", cfg.handlerRevokeRefreshToken)

	serveMux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	serveMux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGetAll)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGetByID)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirpByID)

	serveMux.HandleFunc("POST /admin/reset", cfg.ResetHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.GetMetrics)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, request)
	})
}
