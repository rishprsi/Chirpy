package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) GetMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	response := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	// status, err := writer.Write([]byte("Hits: " + string(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load()))))
	status, err := writer.Write([]byte(response))
	if err != nil {
		log.Printf("Failed to get the metrics with the following error %s, %v", err, status)
	}
}

func (cfg *apiConfig) ResetMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	cfg.fileserverHits.Store(0)
}
