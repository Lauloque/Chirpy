/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding body: %s", err))
		return
	}

	const maxChirpLength = 140
	if len(params.Body) <= maxChirpLength {
		payload := returnVal{Valid: true}
		respondWithJSON(w, http.StatusBadRequest, payload)
	} else {
		respondWithError(w, http.StatusOK, "Chirp is too long")
	}
}
