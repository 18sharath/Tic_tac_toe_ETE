package handlers

import (
	"encoding/json"
	"net/http"

	"tic_tac_toe/game"
	"tic_tac_toe/store"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateGameRequest struct {
	Mode       game.Mode       `json:"mode"`
	Difficulty game.Difficulty `json:"difficulty"`
	BoardSize  int             `json:"boardSize"`
}

type MoveRequest struct{
	Player string `json:"player"`
	Row int `json:"row"`
	Col int `json:"col"`
}


type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

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
		oMover = game.NewBotMover(req.Difficulty)

	case game.ModeBotVsBot:
		xMover = game.NewBotMover(req.Difficulty)
		oMover = game.NewBotMover(req.Difficulty)

	default:
		http.Error(w, "invalid move", http.StatusBadRequest)
		return
	}

	g := game.NewGame(id, req.BoardSize, req.Mode, req.Difficulty, xMover, oMover)

	if req.Mode == game.ModeBotVsBot {
		for !g.Draw && g.Winner == "" {
			if err := g.Maketurn(); err != nil {
				break
			}
			g.Evaluate()
		}
	}
	h.store.Create(g)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

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

func (h *Handler) MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	g, ok := h.store.Get(id)

	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if g.Winner != "" || g.Draw {
		http.Error(w, "game already finished", http.StatusBadRequest)
		return
	}

	var req MoveRequest

	if err:= json.NewDecoder(r.Body).Decode(&req);err!=nil{
		http.Error(w,"invalid request",http.StatusBadRequest)
		return
	}


	if req.Player != "X" && req.Player != "O" {
		http.Error(w, "invalid player", http.StatusBadRequest)
		return
	}

	var err error

	if g.Turn == "X" && g.PlayerX==nil{
		err=g.MakeMove(req.Player,req.Row,req.Col)
	} else if g.Turn=="O" && g.PlayerO==nil{
		err= g.MakeMove(req.Player,req.Row,req.Col)
	} else {
		err = g.Maketurn()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.Evaluate()

	if g.Winner=="" && !g.Draw{
		if (g.Turn == "X" && g.PlayerX!=nil) || (g.Turn =="O" && g.PlayerO!=nil) {
			if err:=g.Maketurn();err==nil{
				g.Evaluate()
			}
		}
	}

	if err := h.store.Create(g); err != nil {
		http.Error(w, "failed to save game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Delete(id); err != nil {
		http.Error(w, "failed to delete game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
