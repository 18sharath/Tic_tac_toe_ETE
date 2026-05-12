// Package main implements the external bot service for Tic Tac Toe.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// newServer builds and returns the configured HTTP server for the given address.
func newServer(addr string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/move", moveHandler)

	return &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func main() {
	port := flag.String("botservice-port", "9090", "botservice port")
	addr := ":" + *port

	srv := newServer(addr)

	log.Printf("bot service listening on %s\n", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
