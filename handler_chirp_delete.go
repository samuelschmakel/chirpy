package main

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/samuelschmakel/chirpy/internal/auth"
)

func (cfg *apiconfig) handlerChirpDelete(w http.ResponseWriter, req *http.Request) {
	aTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find token", err)
		return
	}

	// Next, validate the token:
	userID, err := auth.ValidateJWT(aTokenString, os.Getenv("SECRET_KEY"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirpIDString := req.PathValue("chirpID")
	if chirpIDString == "" {
		respondWithError(w, http.StatusBadRequest, "ID is required", nil)
		return
	}

	chirpUUID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID", nil)
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found in database", err)
		return
	}

	if dbChirp.UserID.String() != userID.String() {
		respondWithError(w, http.StatusForbidden, "can't delete another user's chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
