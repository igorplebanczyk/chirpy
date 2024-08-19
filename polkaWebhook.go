package main

import (
	"encoding/json"
	"net/http"
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
