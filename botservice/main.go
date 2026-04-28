// Package main implements the external bot service for Tic Tac Toe.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("port", "9090", "botservice port")
	flag.Parse()

	addr := ":" + *port

	r := mux.NewRouter()

	r.HandleFunc("/move", moveHandler)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("bot service listening on %s\n", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
