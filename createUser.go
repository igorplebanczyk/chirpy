package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	fmt.Println("creating a decoder")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	fmt.Println("decoding finished")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if !strings.Contains(params.Email, "@") || !strings.Contains(params.Email, ".") {
		respondWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}
	fmt.Println("creating a user")
	user, err := cfg.db.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		ID:    user.ID,
		Email: user.Email,
	})
	fmt.Println("user created")
}
