package game

import (
	"errors"
)

// Game represents the state of the single game instance.
type Game struct {
	ID         string     `json:"id"`
	Board      Board      `json:"board"`
	PlayerX    Mover      `json:"-"`
	PlayerO    Mover      `json:"-"`
	Turn       string     `json:"turn"`
	Winner     string     `json:"winner"`
	Draw       bool       `json:"draw"`
	Mode       Mode       `json:"mode"`
	Difficulty Difficulty `json:"difficulty"`
}

// NewGame initialze the game with provided data
func NewGame(id string, size int, mode Mode, difficulty Difficulty, xMover Mover, oMover Mover) *Game {
	return &Game{
		ID:         id,
		Board:      NewBoard(size),
		PlayerX:    xMover,
		PlayerO:    oMover,
		Turn:       "X",
		Mode:       mode,
		Difficulty: difficulty,
	}
}

// PlayTurn processes a move for both human and bot players.
// It handles human input, triggers a bot move if required,
// and evaluates the game state after each move.
func (g *Game) PlayTurn(player string, row, col int) error{
	if g.Winner != "" || g.Draw {
		return errors.New("game already finished")
	}


	if player != g.Turn {
		return errors.New("not your turn")
	}

	
	if (g.Turn == "X" && g.PlayerX == nil) || (g.Turn == "O" && g.PlayerO == nil) {
		if err := g.MakeMove(player, row, col); err != nil {
			return err
		}

		g.Evaluate()
	} else {
			if err := g.Maketurn(); err != nil {
			return err
		}

		g.Evaluate()	
	}

	if g.Winner != "" || g.Draw {
		return nil
	}

	if (g.Turn == "X" && g.PlayerX != nil) || (g.Turn == "O" && g.PlayerO != nil) {
		if err := g.Maketurn(); err != nil {
			return err
		}

		g.Evaluate()
	}

	return nil
}


// Maketurn helps to place bot player move in game board
func (g *Game) Maketurn() error {
	var mover Mover
	if g.Turn == "X" {
		mover = g.PlayerX
	} else {
		mover = g.PlayerO
	}

	if mover == nil {
		return errors.New("no mover available for current player")
	}

	pos, err := mover.Move(g.Board, g.Turn)
	if err != nil {
		return err
	}
	if g.Board[pos.Row][pos.Col] != "" {
		return errors.New("invalid move")
	}
	g.Board[pos.Row][pos.Col] = g.Turn
	g.toggleTurn()
	return nil
}

// Evaluate checks the current board state and updates the winner or draw status.
func (g *Game) Evaluate() {
	if winner, ok := g.checkRows(); ok {
		g.Winner = winner
		return
	}

	if winner, ok := g.checkCols(); ok {
		g.Winner = winner
		return
	}

	if winner, ok := g.checkMainDiagonal(); ok {
		g.Winner = winner
		return
	}

	if winner, ok := g.checkAntiDiagonal(); ok {
		g.Winner = winner
		return
	}

	if g.checkDraw() {
		g.Draw = true
	}
}

// checkRows verifies if any row has identical non-empty values.
func (g *Game) checkRows() (string, bool) {
	size := len(g.Board)

	for i := 0; i < size; i++ {
		first := g.Board[i][0]
		if first == "" {
			continue
		}

		win := true
		for j := 1; j < size; j++ {
			if g.Board[i][j] != first {
				win = false
				break
			}
		}

		if win {
			return first, true
		}
	}

	return "", false
}

// checkCols verifies if any column has identical non-empty values.
func (g *Game) checkCols() (string, bool) {
	size := len(g.Board)

	for j := 0; j < size; j++ {
		first := g.Board[0][j]
		if first == "" {
			continue
		}

		win := true
		for i := 1; i < size; i++ {
			if g.Board[i][j] != first {
				win = false
				break
			}
		}

		if win {
			return first, true
		}
	}

	return "", false
}

// checkMainDiagonal verifies the primary diagonal for a winning condition.
func (g *Game) checkMainDiagonal() (string, bool) {
	size := len(g.Board)

	first := g.Board[0][0]
	if first == "" {
		return "", false
	}

	for i := 1; i < size; i++ {
		if g.Board[i][i] != first {
			return "", false
		}
	}

	return first, true
}

// checkAntiDiagonal verifies the secondary diagonal for a winning condition.
func (g *Game) checkAntiDiagonal() (string, bool) {
	size := len(g.Board)

	first := g.Board[0][size-1]
	if first == "" {
		return "", false
	}

	for i := 1; i < size; i++ {
		if g.Board[i][size-1-i] != first {
			return "", false
		}
	}

	return first, true
}

// checkDraw determines if the board is full with no winner.
func (g *Game) checkDraw() bool {
	size := len(g.Board)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if g.Board[i][j] == "" {
				return false
			}
		}
	}

	return true
}

// MakeMove helps to place human player move in game board
func (g *Game) MakeMove(player string, row, col int) error {
	if player != g.Turn {
		return errors.New("not your turn")
	}
	size := len(g.Board)

	if row < 0 || row >= size || col < 0 || col >= size {
		return errors.New("invalid position")
	}

	if g.Board[row][col] != "" {
		return errors.New("cell already occupied")
	}

	g.Board[row][col] = player

	g.toggleTurn()
	return nil
}

// toggleTurn toggle player on each move
func (g *Game) toggleTurn() {
	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}
}
