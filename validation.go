package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ChirpValidation(writer http.ResponseWriter, request *http.Request) {
	type validateBody struct {
		Body string `json:"body"`
	}

	type validateRespBody struct {
		CleanedBody string `json:"cleaned_body"`
		Error       string `json:"error"`
		Valid       bool   `json:"valid"`
	}
	reqBody := validateBody{}
	respBody := validateRespBody{}
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Failed to decode incoming body for chirp validation with err: %v", err)
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(400)
		respBody.Error = "Something went wrong"
		body, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Failed to encode response body for chirp validation with err: %v", err)
		}
		status, err := writer.Write(body)
		if err != nil {
			log.Printf("Failed to write the final response body %v %v", err, status)
		}
		return
	}

	statusCode := 200
	if len(reqBody.Body) > 140 {
		log.Printf("Chirp too long for the req body %s", reqBody.Body)
		respBody.Error = "Chirp is too long"
		statusCode = 400
	} else {
		respBody.Valid = true
		cleanedBody := reqBody.Body
		badWords := [3]string{"kerfuffle", "sharbert", "fornax"}
		for _, badWord := range badWords {
			words := strings.Split(cleanedBody, " ")
			fmt.Printf("words are: %v", words)
			for index, word := range words {
				if badWord == strings.ToLower(word) {
					words[index] = "****"
				}
			}
			cleanedBody = strings.Join(words, " ")
		}
		respBody.CleanedBody = cleanedBody
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	body, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Failed to decode incoming body for chirp validation with err: %v", err)
	}
	status, err := writer.Write(body)
	if err != nil {
		log.Printf("Failed to write the final response body %v %v", err, status)
	}
}
