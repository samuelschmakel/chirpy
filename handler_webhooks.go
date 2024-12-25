package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/samuelschmakel/chirpy/internal/auth"
)

func (cfg *apiconfig) handlerWebhooks(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no API key in request header", err)
		return
	}

	if apiKey != cfg.polkakey {
		respondWithError(w, http.StatusUnauthorized, "invalid key", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid userID", err)
		return
	}
	_, err = cfg.db.UpgradeToChirpyRed(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't update user to chirpy red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
