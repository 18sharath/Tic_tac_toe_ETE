// Package main starts the Tic Tac Toe backend server.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"tic_tac_toe/game"
	"tic_tac_toe/handlers"
	"tic_tac_toe/store"

	"github.com/gorilla/mux"
)

// corsMiddleware adds CORS headers to allow browser clients to access the API.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	storeType := flag.String("store", "memory", "memory or file")
	port := flag.String("port", "8080", "server port")
	botserviceURL := flag.String("bot-service-url", "http://localhost:9090/move", "bot service endpoint")

	flag.Parse()

	game.SetBotServiceURL(*botserviceURL)

	var s store.GameStore

	if *storeType == "file" {
		s = store.NewFileStore("data")
	} else {
		s = store.NewMemoryStore()
	}

	r := mux.NewRouter()
	handler := handlers.NewHandler(s)

	// API routes
	r.HandleFunc("/games", handler.CreateGameHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/games/{id}", handler.GetGameHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/games/{id}", handler.MakeMoveHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/games/{id}", handler.DeleteGameHandler).Methods("DELETE", "OPTIONS")

	addr := ":" + *port

	log.Printf("Server running on %v", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(r),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
