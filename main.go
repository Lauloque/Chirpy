/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Lauloque/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polka_key      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       os.Getenv("PLATFORM"),
		secret:         os.Getenv("SECRET"),
		polka_key:      os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))),
	)

	// USERS
	mux.HandleFunc("POST /api/users", apiCfg.handleUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUsersUpdate)
	// Authentification
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	// Chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handleChirpsList)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleChirpsDelete)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirpsCreate)
	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)
	// Admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
