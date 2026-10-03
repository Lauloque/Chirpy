/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding body: %s", err))
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cleaned := cleanupBody(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVal{Cleaned_body: cleaned})
}

func cleanupBody(body string, badWords []string) string {
	words := strings.Split(body, " ")
	for i := range words {
		for _, badWord := range badWords {
			if strings.ToLower(words[i]) == badWord {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
