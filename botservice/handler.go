package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func moveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateBoard(req.Board); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player := req.Player
	if player == "" {
		player = "O"
	}

	move := bestMove(req.Board, player)
	if move.Row == -1 {
		http.Error(w, "no moves available", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(move); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func validateBoard(board [][]string) error {
	if len(board) == 0 {
		return errors.New("board must not be empty")
	}
	size := len(board)
	for _, row := range board {
		if len(row) != size {
			return errors.New("board must be square (NxN)")
		}
	}
	return nil
}
