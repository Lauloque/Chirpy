/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"net/http"

	"github.com/Lauloque/Chirpy/internal/auth"
	"github.com/Lauloque/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type repsponse struct {
		User
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	query_params := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hash,
	}
	user, err := cfg.db.UpdateUser(r.Context(), query_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update credentials", err)
		return
	}

	respondWithJSON(w, http.StatusOK, repsponse{User: User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}})
}
