package main

import (
	"encoding/json"
	"net/http"

	"github.com/samuelschmakel/chirpy/internal/auth"
	"github.com/samuelschmakel/chirpy/internal/database"
)

func (cfg *apiconfig) handlerUsersCreate(w http.ResponseWriter, req *http.Request) {
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

	dbUser, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
	}

	passParams := database.AddPasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
	}

	_, err = cfg.db.AddPassword(req.Context(), passParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't store password", err)
	}

	respondWithJSON(w, http.StatusCreated, user)
}
