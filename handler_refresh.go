/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"net/http"
	"time"

	"github.com/Lauloque/Chirpy/internal/auth"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid refresh token in the header", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token user, or user's token is expired or revoked.", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
