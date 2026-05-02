// Package game contains the core Tic Tac Toe game logic
package game

// Position represents location on the  grame board.
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Board represents the game board used to play game.
type Board [][]string

// NewBoard creates a two dimensional game board based on size
func NewBoard(size int) Board {
	board := make(Board, size)
	for i := range board {
		board[i] = make([]string, size)
	}
	return board
}

// ApplyMove returns a new Board with the given move applied.
func ApplyMove(b Board, pos Position, player string) Board {
	size := len(b)
	if size == 0 {
		return b
	}

	nb := make(Board, size)
	for i := range b {
		nb[i] = make([]string, len(b[i]))
		copy(nb[i], b[i])
	}

	if pos.Row >= 0 && pos.Row < size && pos.Col >= 0 && pos.Col < len(nb[pos.Row]) {
		nb[pos.Row][pos.Col] = player
	}

	return nb
}
