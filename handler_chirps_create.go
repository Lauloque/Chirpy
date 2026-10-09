/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Lauloque/Chirpy/internal/auth"
	"github.com/Lauloque/Chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode chirp request", err)
		return
	}

	// USER VALIDATION
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	authUserId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT:", err)
		return
	}
	// END USER VALIDATION

	// Chirp VALIDATION

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cleaned := cleanupBody(params.Body, badWords)

	chirpParams := database.ChirpCreateParams{
		Body:   cleaned,
		UserID: authUserId,
	}

	// END VALIDATION

	chirp, err := cfg.db.ChirpCreate(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp: %v", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
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
