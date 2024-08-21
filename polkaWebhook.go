package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	polkaAPIKey, err := retrievePolkaAPIKeyFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if polkaAPIKey != cfg.polkaAPIKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	err = cfg.db.UpdateUserChirpyRedStatus(params.Data.UserID, true)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func retrievePolkaAPIKeyFromHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("no api key provided")
	}

	if !strings.HasPrefix(token, "ApiKey ") {
		return "", errors.New("invalid api key format")
	}

	token = strings.TrimPrefix(token, "ApiKey ") // Trim the "Bearer " prefix to get the actual token
	if token == "" {
		return "", errors.New("malformed api key")
	}

	return token, nil
}
