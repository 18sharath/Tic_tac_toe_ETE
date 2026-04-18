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
