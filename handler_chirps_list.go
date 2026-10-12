/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"net/http"

	"github.com/Lauloque/Chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleChirpsList(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	chirps := []database.Chirp{}

	if authorID != uuid.Nil {
		chirps, err = cfg.db.ChirpsListFromAuthorId(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.ChirpsList(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't list chirps: %v", err)
		return
	}

	chirpsArray := []Chirp{}
	for i := range chirps {
		chirpsArray = append(chirpsArray, Chirp{
			Id:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserId:    chirps[i].UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsArray)
}

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}

	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, nil
	}
	return authorID, nil
}
