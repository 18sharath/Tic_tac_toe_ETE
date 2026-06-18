// Package main provides a simple static file server for the frontend.
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "3000", "frontend server port")
	flag.Parse()

	fs := http.FileServer(http.Dir("."))

	log.Printf("Frontend server running on http://localhost:%s", *port)
	log.Printf("Make sure backend is running on http://localhost:8080")

	err := http.ListenAndServe(":"+*port, fs)
	if err != nil {
		log.Fatal(err)
	}
}
