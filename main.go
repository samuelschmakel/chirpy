package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	fileSystem := http.Dir(filepathRoot)

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(fileSystem)))
	mux.HandleFunc("/healthz", handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK) // status code is 200
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
