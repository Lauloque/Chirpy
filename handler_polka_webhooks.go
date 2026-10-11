/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"net/http"

	"github.com/Lauloque/Chirpy/internal/auth"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	type response struct {
		Body string
	}

	apikey, err := auth.GetAPIKey(r.Header)

	if err != nil || apikey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key", err)
		return
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade user:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
