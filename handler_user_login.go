package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/samuelschmakel/chirpy/internal/auth"
)

func (cfg *apiconfig) handlerUserLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"` // Optional field (pointer)
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

	fmt.Printf("Retrieved password: %s\n", dbUser.HashedPassword)
	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Making a JWT for the user:
	expiresIn := 3600 * time.Second
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds < 3600 {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.secretkey, expiresIn)
	if err != nil {
		fmt.Println("couldn't make jwt")
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     tokenString,
	}

	respondWithJSON(w, http.StatusOK, user)

}
