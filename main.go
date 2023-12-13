package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// 1. Add the readiness endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Write the Content-Type header
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// Write the status code
		w.WriteHeader(http.StatusOK)
		// Write the body text
		w.Write([]byte("OK"))
	})
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir("logo.png")))

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
