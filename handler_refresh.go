package main

import (
	"net/http"
	"os"

	"github.com/samuelschmakel/chirpy/internal/auth"
)

func (cfg *apiconfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find token in header", err)
		return
	}

	// Next, validate the token:
	userID, err := auth.ValidateJWT(tokenString, os.Getenv("SECRET_KEY"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	rToken, err := cfg.db.GetTokenFromUser(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find token from database", err)
		return
	}

	if rToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token is expired", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user from refresh token", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretkey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create an access token", err)
		return
	}

	t := tokenStruct{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, t)
}
