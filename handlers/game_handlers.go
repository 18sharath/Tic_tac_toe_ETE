// Package handlers contains HTTP handlers for managing Tic Tac Toe games,
// including creating games, making moves, retrieving game state,
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tic_tac_toe/game"
	"tic_tac_toe/store"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateGameRequest represents the playload required to create a new game.
type CreateGameRequest struct {
	Mode        game.Mode       `json:"mode"`
	DifficultyX game.Difficulty `json:"difficultyX"`
	DifficultyO game.Difficulty `json:"difficultyO"`
	BoardSize   int             `json:"boardSize"`
}

// MoveRequest represents the payload required to make a move in a game.
type MoveRequest struct {
	Player string `json:"player"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

// Handler handles http request using gamestore
type Handler struct {
	store store.GameStore
}

// NewHandler creates a new handler with the given gamestore
func NewHandler(s store.GameStore) *Handler {
	return &Handler{store: s}
}

// CreateGameHandler handles the http request for create a new game
func (h *Handler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.BoardSize < 0 {
		req.BoardSize = 3
	}

	id := uuid.New().String()

	var xMover game.Mover
	var oMover game.Mover

	switch req.Mode {
	case game.ModeHumanVsHuman:
		xMover = nil
		oMover = nil

	case game.ModeHumanVsBot:
		xMover = nil
		oMover = game.NewBotMover(req.DifficultyO)

	case game.ModeBotVsBot:
		xMover = game.NewBotMover(req.DifficultyX)
		oMover = game.NewBotMover(req.DifficultyO)

	default:
		http.Error(w, "invalid move", http.StatusBadRequest)
		return
	}

	g := game.NewGame(id, req.BoardSize, req.Mode, req.DifficultyO, xMover, oMover)

	if req.Mode == game.ModeBotVsBot {
		for !g.Draw && g.Winner == "" {
			if err := g.Maketurn(); err != nil {
				break
			}
			g.Evaluate()
		}
	}
	if err := h.store.Create(g); err != nil {
		log.Println("error creating game:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetGameHandler handles the http request for get games based on gameId
func (h *Handler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	g, ok := h.store.Get(id)

	if !ok {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// getGameFromRequest retrieves game from store using request ID.
func (h *Handler) getGameFromRequest(r *http.Request) (*game.Game, error) {
	id := mux.Vars(r)["id"]

	g, ok := h.store.Get(id)
	if !ok {
		return nil, fmt.Errorf("game not found")
	}

	return g, nil
}

// validateGameNotFinished ensures game is still active.
func validateGameNotFinished(g *game.Game) error {
	if g.Winner != "" || g.Draw {
		return fmt.Errorf("game already finished")
	}
	return nil
}

// decodeMoveRequest parses the incoming move request body.
func decodeMoveRequest(r *http.Request) (MoveRequest, error) {
	var req MoveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, fmt.Errorf("invalid request")
	}

	return req, nil
}

// validatePlayer checks if player is valid (X or O).
func validatePlayer(player string) error {
	if player != "X" && player != "O" {
		return fmt.Errorf("invalid player")
	}
	return nil
}

// writeJSONResponse writes JSON response to client.
func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err:=json.NewEncoder(w).Encode(data);err!=nil{
		http.Error(w,"failed to encode response", http.StatusInternalServerError)
		return
	}
}

// MakeMoveHandler handles the http request to make move.
func (h *Handler) MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	g, err := h.getGameFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := validateGameNotFinished(g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := decodeMoveRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validatePlayer(req.Player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := g.PlayTurn(req.Player, req.Row, req.Col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Create(g); err != nil {
		http.Error(w, "failed to save game", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, g)
}

// DeleteGameHandler hanldes the http request to delete already existing game
func (h *Handler) DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Delete(id); err != nil {
		http.Error(w, "failed to delete game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
