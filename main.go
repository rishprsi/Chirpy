package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/rishprsi/Chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to create a connection with the database with err: %v\n", err)
	}
	dbqueries := database.New(db)
	port := "8080"
	filepathRoot := "."
	serveMux := http.NewServeMux()
	cfg := apiConfig{}

	serveMux.Handle("/app", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	serveMux.HandleFunc("GET /api/healthz", Readiness)
	serveMux.HandleFunc("POST /api/validate_chirp", ChirpValidation)

	serveMux.HandleFunc("POST /admin/reset", cfg.ResetMetrics)
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
