// Package main starts the Tic Tac Toe backend server.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"tic_tac_toe/game"
	"tic_tac_toe/handlers"
	"tic_tac_toe/store"

	"github.com/gorilla/mux"
)

// BotServiceConfig holds the configuration for available bot services.
type BotServiceConfig struct {
	Services []Botservice `json:"services"`
}

// Botservice represents a single bot service with its connection details.
type Botservice struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

var botServiceConfig BotServiceConfig

func loadBotServiceConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		botServiceConfig = BotServiceConfig{Services: []Botservice{}} // artha agilla
		return nil
	}
	return json.Unmarshal(data, &botServiceConfig)
}

func getBotServicesHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(botServiceConfig); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

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
	botConfigPath := flag.String("bot-config", "config/bot_services.json", "path to bot services config file")

	flag.Parse()

	if err := loadBotServiceConfig(*botConfigPath); err != nil {
		log.Printf("Warning: could not load bot config: %v", err)
	}

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
	r.HandleFunc("/bot-services", getBotServicesHandler).Methods("GET", "OPTIONS")

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
