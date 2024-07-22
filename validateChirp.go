package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, validResponse{CleanedBody: replaceBadWords(params.Body)})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(errorResponse{Error: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func replaceBadWords(body string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	replacement := "****"

	words := strings.Fields(body)
	for _, badWord := range badWords {
		for i, word := range words {
			if strings.ToLower(word) == badWord {
				words[i] = replacement
			}
		}
	}
	body = strings.Join(words, " ")

	return body
}
