package main

import (
	"go-server/internal/database"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	queryValue := r.URL.Query().Get("author_id")

	if queryValue != "" {
		authorID, err := strconv.Atoi(queryValue)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}

		chirps, err := cfg.db.GetChirpsByAuthorID(authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Author not found")
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	}

	var chirps []database.Chirp
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	queryValue = r.URL.Query().Get("sort")

	if queryValue == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return dbChirps[i].ID > dbChirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpHandlerSpecific(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	chirp, err := cfg.db.GetChirpByID(requestedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
