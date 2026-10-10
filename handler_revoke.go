/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"net/http"

	"github.com/Lauloque/Chirpy/internal/auth"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid refresh token in the header", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
