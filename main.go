package main

import (
	"log"
	"net/http"

	"tic_tac_toe/handlers"
	"tic_tac_toe/store"
)

func main() {
    err:=store.LoadGames()
    if err!=nil{
        log.Fatal("Failed")
    }
    http.HandleFunc("/games", handlers.CreateGameHandler)
    http.HandleFunc("/games/", handlers.GameHandler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}