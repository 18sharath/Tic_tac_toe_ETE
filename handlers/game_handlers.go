package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/google/uuid"
    "tic_tac_toe/game"
    "tic_tac_toe/store"
)

type CreateGameRequest struct {
    Mode       int `json:"mode"`
    Difficulty int `json:"difficulty"`
}

type MoveRequest struct {
    Row int `json:"row"`
    Col int `json:"col"`
}


func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var req CreateGameRequest
    json.NewDecoder(r.Body).Decode(&req)

    id := uuid.New().String()
    g := game.NewGame(id, req.Mode, req.Difficulty)

	if g.Mode == 3 {
    	for {
			g.BotMove()
			winner := g.CheckWinner();
			if winner != "" {
				g.Winner = winner
				break
			}

			if g.CheckDraw() {
				g.Draw = true
				break
			}

			g.TogglePlayer()
		}
  	}
    store.Mutex.Lock()
    store.Games[id] = g
    store.Mutex.Unlock()
    store.SaveGame()
    json.NewEncoder(w).Encode(g)
}


func GameHandler(w http.ResponseWriter, r *http.Request) {
    
    id := strings.Split(r.URL.Path, "/")[2]
    g, ok := store.Games[id]

    if !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    switch r.Method {

    case http.MethodGet:
        json.NewEncoder(w).Encode(g)

    case http.MethodPost:
        var move MoveRequest
        json.NewDecoder(r.Body).Decode(&move)

        
		err := g.MakeMove(move.Row, move.Col)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

        if winner := g.CheckWinner(); winner != "" {
            g.Winner = winner
            json.NewEncoder(w).Encode(g)
            return
        }

        if g.CheckDraw() {
            g.Draw = true
            json.NewEncoder(w).Encode(g)
            return
        }

       g.TogglePlayer()

       // Player vs Bot
        if g.Mode == 2 && g.Player == "O" {
            g.BotMove()
            g.TogglePlayer()
        }

        if winner := g.CheckWinner(); winner != "" {
            g.Winner = winner
        }

        if g.CheckDraw() {
            g.Draw = true
        }
        store.SaveGame()
        json.NewEncoder(w).Encode(g)
    }
}