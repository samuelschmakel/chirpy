package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/samuelschmakel/chirpy/internal/database"
)

func (cfg *apiconfig) handlerChirpsRetrieve(w http.ResponseWriter, req *http.Request) {
	authorIDString := req.URL.Query().Get("author_id")
	sort := req.URL.Query().Get("sort")

	if sort != "asc" && sort != "desc" && sort != "" {
		respondWithError(w, http.StatusBadRequest, "invalid value for sort query parameter", fmt.Errorf("invalid query"))
		return
	}

	var authorID uuid.UUID
	var err error

	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid user id", err)
			return
		}
	} else {
		authorID, _ = uuid.Parse("")
	}

	var chirps []database.Chirp

	if authorIDString != "" {
		if sort != "desc" {
			chirps, err = cfg.db.GetChirpsFromUser(req.Context(), authorID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid request", err)
				return
			}
		} else {
			chirps, err = cfg.db.GetChirpsFromUserDesc(req.Context(), authorID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid request", err)
				return
			}
		}
	} else {
		if sort != "desc" {
			chirps, err = cfg.db.GetChirps(req.Context())
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid request", err)
				return
			}
		} else {
			chirps, err = cfg.db.GetChirpsDesc(req.Context())
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid request", err)
				return
			}
		}
	}

	chirpArr := []Chirp{}

	for _, chirp := range chirps {
		chirpArr = append(chirpArr, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpArr)
}

func (cfg *apiconfig) handlerChirpRetrieve(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID is required", nil)
		return
	}

	u, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID", nil)
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), u)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "id not found", err)
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}
