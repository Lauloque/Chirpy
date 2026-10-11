/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"net/http"

	"github.com/Lauloque/Chirpy/internal/auth"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	// USER VALIDATION
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT:", err)
		return
	}

	// END USER VALIDATION

	// CHIRP VALIDATION

	chirp, err := cfg.db.ChirpsGet(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp:", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User doesn't own this chirp:", nil)
		return
	}

	// END CHIRP VALIDATION

	err = cfg.db.ChirpDelete(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
