// Package main implements the external bot service for Tic Tac Toe.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"flag"
	"time"
)


type moveRequest struct{
	Board [][] string `json:"board"`
	Player string	`json:"player"`
}

type moveResponse struct{
	Row int `json:"row"`
	Col int `json:"col"`
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	var req moveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	board := req.Board
	bestScore := -1000
	bestMove := moveResponse{
		Row: -1,
		Col: -1,
	}

	// Bot always plays as O
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == "" {
				board[i][j] = "O"

				score := minimax(board, false)

				board[i][j] = ""

				if score > bestScore {
					bestScore = score
					bestMove = moveResponse{
						Row: i,
						Col: j,
					}
				}
			}
		}
	}

	if bestMove.Row == -1 {
		http.Error(w, "no moves available", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(bestMove); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func checkWinner(board [][]string) string {
	size := len(board)

	// rows
	for i := 0; i < size; i++ {
		if board[i][0] != "" {
			win := true
			for j := 1; j < size; j++ {
				if board[i][j] != board[i][0] {
					win = false
					break
				}
			}
			if win {
				return board[i][0]
			}
		}
	}

	// cols
	for j := 0; j < size; j++ {
		if board[0][j] != "" {
			win := true
			for i := 1; i < size; i++ {
				if board[i][j] != board[0][j] {
					win = false
					break
				}
			}
			if win {
				return board[0][j]
			}
		}
	}

	// main diagonal
	if board[0][0] != "" {
		win := true
		for i := 1; i < size; i++ {
			if board[i][i] != board[0][0] {
				win = false
				break
			}
		}
		if win {
			return board[0][0]
		}
	}

	// anti diagonal
	if board[0][size-1] != "" {
		win := true
		for i := 1; i < size; i++ {
			if board[i][size-1-i] != board[0][size-1] {
				win = false
				break
			}
		}
		if win {
			return board[0][size-1]
		}
	}

	return ""
}

func isBoardFull(board [][]string) bool {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func minimax(board [][]string, isMaximizing bool) int {
	winner := checkWinner(board)

	if winner == "O" {
		return 10
	}

	if winner == "X" {
		return -10
	}

	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		bestScore := -1000

		for i := 0; i < len(board); i++ {
			for j := 0; j < len(board[i]); j++ {
				if board[i][j] == "" {
					board[i][j] = "O"

					score := minimax(board, false)

					board[i][j] = ""

					if score > bestScore {
						bestScore = score
					}
				}
			}
		}

		return bestScore
	}

	bestScore := 1000

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == "" {
				board[i][j] = "X"

				score := minimax(board, true)

				board[i][j] = ""

				if score < bestScore {
					bestScore = score
				}
			}
		}
	}

	return bestScore
}


func main(){
	port := flag.String("port", "9090", "botservice port")
	flag.Parse()

	addr := ":" + *port

	http.HandleFunc("/move",moveHandler)
	
	log.Printf("Bot Service running on %v",addr)

	srv:= &http.Server{
		Addr: addr,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
		IdleTimeout: 120*time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}