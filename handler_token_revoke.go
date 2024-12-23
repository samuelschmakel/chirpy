package main

import (
	"net/http"

	"github.com/samuelschmakel/chirpy/internal/auth"
)

func (cfg *apiconfig) handlerTokenRevoke(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find token", err)
		return
	}

	_, err = cfg.db.RevokeToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
