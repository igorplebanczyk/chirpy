package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	userID, err := cfg.GetAuthenticatedUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate user")
		return
	}

	chirp, err := cfg.db.GetChirpByID(requestedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.db.DeleteChirp(requestedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
