package main

import (
	"flag"
	"log"
	"net/http"

	"tic_tac_toe/handlers"
	"tic_tac_toe/store"

	"github.com/gorilla/mux"
)

func main() {
	storeType := flag.String("store", "memory", "memory or file")
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	var s store.GameStore

	if *storeType == "file" {
		s = store.NewFileStore("data")
	} else {
		s = store.NewMemoryStore()
	}

	r := mux.NewRouter()
	handler := handlers.NewHandler(s)

	r.HandleFunc("/games", handler.CreateGameHandler).Methods("POST")
	r.HandleFunc("/games/{id}", handler.GetGameHandler).Methods("GET")
	r.HandleFunc("/games/{id}", handler.MakeMoveHandler).Methods("PUT")
	r.HandleFunc("/games/{id}", handler.DeleteGameHandler).Methods("DELETE")

	addr := ":" + *port

	log.Printf("Server running on %v", addr)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}

}
