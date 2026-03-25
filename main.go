package main

import (
    "log"
    "net/http"

    "tic_tac_toe/handlers"
)

func main() {
    http.HandleFunc("/games", handlers.CreateGameHandler)
    http.HandleFunc("/games/", handlers.GameHandler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}