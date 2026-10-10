/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lauloque/Chirpy/internal/auth"
	"github.com/Lauloque/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password.", err)
		return
	}

	success, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password check failed:", err)
		return
	}
	if !success {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password.", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT:", err)
		return
	}

	rereshTokenString := auth.MakeRefreshToken()
	rtParams := database.CreateRefreshTokenParams{
		Token:     rereshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), rtParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token:", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}
