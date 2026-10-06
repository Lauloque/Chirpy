/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"net/http"

	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleChirpsList(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps: %v", err)
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
