// Package handlers contains HTTP handlers for managing Tic Tac Toe games,
// including creating games, making moves, retrieving game state,
package handlers

import (
	"encoding/json"
<<<<<<< HEAD
	"fmt"
=======
	"log"
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
	"net/http"
	"tic_tac_toe/game"
	"tic_tac_toe/store"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateGameRequest represents the playload required to create a new game.
type CreateGameRequest struct {
<<<<<<< HEAD
<<<<<<< HEAD
	Mode        game.Mode       `json:"mode"`
	DifficultyX game.Difficulty `json:"difficultyX"`
	DifficultyO game.Difficulty `json:"difficultyO"`
	BoardSize   int             `json:"boardSize"`
=======
	Mode       game.Mode       `json:"mode"`
	DifficultyX game.Difficulty `json:"difficultyX"`
	DifficultyO game.Difficulty `json:"difficultyO"`
	BoardSize  int             `json:"boardSize"`
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
	Mode        game.Mode       `json:"mode"`
	DifficultyX game.Difficulty `json:"difficultyX"`
	DifficultyO game.Difficulty `json:"difficultyO"`
	BoardSize   int             `json:"boardSize"`
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
}

// MoveRequest represents the payload required to make a move in a game.
type MoveRequest struct {
	Player string `json:"player"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// Handler handles http request using gamestore
=======

>>>>>>> 5f61e4f (This change will ask name of the player and dynamic board size in cli)
=======
// handler handles http request using gamestore
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
// Handler handles http request using gamestore
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
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

	if req.BoardSize < 3 {
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
		http.Error(w, "invalid mode", http.StatusBadRequest)
		return
	}

	g := game.NewGame(id, req.BoardSize, req.Mode, req.DifficultyO, xMover, oMover)

	if req.Mode == game.ModeBotVsBot {
		runBotGame(g)
	}
<<<<<<< HEAD
	if !h.saveGame(w, g) {
=======
	if err := h.store.Create(g); err != nil {
		log.Println("error creating game:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
		return
	}
	writeJSONResponse(w, http.StatusCreated, g)
}

// GetGameHandler handles the http request for get games based on gameId
func (h *Handler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	g, ok := h.store.Get(id)

	if !ok {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}
	writeJSONResponse(w, http.StatusOK, g)
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

// writeJSONResponse writes a JSON response with the given status code.
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
// runBotGame plays a BotVsBot game to completion, stopping on mover error.
func runBotGame(g *game.Game) {
	for !g.Draw && g.Winner == "" {
		if err := g.Maketurn(); err != nil {
			break
		}
		g.Evaluate()
	}
}

// saveGame persists the game and writes a 500 error on failure.
func (h *Handler) saveGame(w http.ResponseWriter, g *game.Game) bool {
	if err := h.store.Create(g); err != nil {
		http.Error(w, "failed to save game", http.StatusInternalServerError)
		return false
	}
	return true
}

// parseMoveRequest decodes and validates the move request body.
func parseMoveRequest(r *http.Request) (MoveRequest, error) {
	req, err := decodeMoveRequest(r)
	if err != nil {
		return req, err
	}
	return req, validatePlayer(req.Player)
}

=======
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
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
=======
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
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
	}
	return nil
}

<<<<<<< HEAD
<<<<<<< HEAD
	req, err := parseMoveRequest(r)
=======
=======
// decodeMoveRequest parses the incoming move request body.
func decodeMoveRequest(r *http.Request) (MoveRequest, error) {
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
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

<<<<<<< HEAD
	var err error

	if g.Turn == "X" && g.PlayerX == nil {
		err = g.MakeMove(req.Player, req.Row, req.Col)
	} else if g.Turn == "O" && g.PlayerO == nil {
		err = g.MakeMove(req.Player, req.Row, req.Col)
	} else {
		err = g.Maketurn()
	}

>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
	req, err := decodeMoveRequest(r)
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
	if err := g.PlayTurn(req.Player, req.Row, req.Col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
=======
	g.Evaluate()
=======
	if err := validatePlayer(req.Player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)

	if err := g.PlayTurn(req.Player, req.Row, req.Col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Create(g); err != nil {
		http.Error(w, "failed to save game", http.StatusInternalServerError)
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
		return
	}

<<<<<<< HEAD
	if !h.saveGame(w, g) {
		return
	}

	writeJSONResponse(w, http.StatusOK, g)
=======
	writeJSONResponse(w, g)
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
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
