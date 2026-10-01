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
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	mux.HandleFunc(readinessPath, readinessHandler)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
