package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if !strings.Contains(params.Email, "@") || !strings.Contains(params.Email, ".") {
		respondWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user, err := cfg.db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
