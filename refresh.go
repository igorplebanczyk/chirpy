package main

import (
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := RetrieveTokenFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, user, err := cfg.db.GetRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if refreshToken.ExpiresAt < time.Now().Unix() {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accessToken, err := cfg.GenerateAccessToken(user)

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
