/* SPDX-License-Identifier: GPL-3.0-or-later */
package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const readinessPath = "/healthz"
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)
	mux.HandleFunc(readinessPath, handlerReadiness)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
