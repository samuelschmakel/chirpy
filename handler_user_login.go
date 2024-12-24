package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samuelschmakel/chirpy/internal/auth"
	"github.com/samuelschmakel/chirpy/internal/database"
)

func (cfg *apiconfig) handlerUserLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Making a JWT for the user:
	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.secretkey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't make jwt", err)
		return
	}

	// Making a refresh token:
	rTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't make refresh token", err)
	}

	// Add the refresh token to the database:
	tokenParams := database.AddTokenParams{
		Token:  rTokenString,
		UserID: dbUser.ID,
	}

	_, err = cfg.db.AddToken(req.Context(), tokenParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't add refresh token to database", err)
	}

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        tokenString,
		RefreshToken: rTokenString,
	}

	fmt.Printf("user data written in response: %v", user)

	respondWithJSON(w, http.StatusOK, user)
}
