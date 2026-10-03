/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding body: %s", err))
		return
	}
	if len(params.Body) <= 140 {
		type returnVal struct {
			Valid bool `json:"valid"`
		}
		payload := returnVal{Valid: true}
		respondWithJSON(w, 200, payload)
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVal struct {
		Error string `json:"error"`
	}
	respBody := returnVal{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
