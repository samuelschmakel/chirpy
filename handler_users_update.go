package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/samuelschmakel/chirpy/internal/auth"
	"github.com/samuelschmakel/chirpy/internal/database"
)

func (cfg *apiconfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {
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

	rToken, err := cfg.db.GetTokenFromUser(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find token from database", err)
		return
	}

	if rToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token is expired", err)
		return
	}

	// These are the user's NEW email and password
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "malformed request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Now, add the new password and email to the database
	passParams := database.UpdatePasswordParams{
		ID:             userID,
		HashedPassword: hashedPassword,
	}

	_, err = cfg.db.UpdatePassword(req.Context(), passParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't store password", err)
		return
	}

	emailParams := database.UpdateEmailParams{
		Email: params.Email,
		ID:    userID,
	}

	_, err = cfg.db.UpdateEmail(req.Context(), emailParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't store email", err)
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(req.Context(), rToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get user", err)
		return
	}

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		RefreshToken: rToken.Token,
		IsChirpyRed:  dbUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)
}
