package game

import (
	"errors"
	"math/rand"
)

// Mover defines an entity capable of making move on the game board
type Mover interface {
	// Move determines and returns the next position for given player based on current game state.
	Move(board Board, player string) (Position, error)
}

// getEmptyCells returns the empty cells in the current game
func getEmptyCells(board Board) []Position {
	size := len(board)
	var cells []Position

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] == "" {
				cells = append(cells, Position{i, j})
			}
		}
	}
	return cells
}

// isWinningMove checks for an move as a winning move or not
func isWinningMove(board Board, r, c int, symbol string) bool {
	size := len(board)

	// place temporarily
	board[r][c] = symbol
	defer func() { board[r][c] = "" }()

	// check row
	win := true
	for j := 0; j < size; j++ {
		if board[r][j] != symbol {
			win = false
			break
		}
	}
	if win {
		return true
	}

	// check column
	win = true
	for i := 0; i < size; i++ {
		if board[i][c] != symbol {
			win = false
			break
		}
	}
	if win {
		return true
	}

	// check main diagonal
	if r == c {
		win = true
		for i := 0; i < size; i++ {
			if board[i][i] != symbol {
				win = false
				break
			}
		}
		if win {
			return true
		}
	}

	// check anti-diagonal
	if r+c == size-1 {
		win = true
		for i := 0; i < size; i++ {
			if board[i][size-1-i] != symbol {
				win = false
				break
			}
		}
		if win {
			return true
		}
	}

	return false
}

// findWinningMove finds out the winning move and returns the position
func findWinningMove(board Board, symbol string) (Position, bool) {
	emptyCells := getEmptyCells(board)

	for _, pos := range emptyCells {
		if isWinningMove(board, pos.Row, pos.Col, symbol) {
			return pos, true
		}
	}
	return Position{}, false
}

// randomMove finds out the empty cells and return the random number inside empty cell
func randomMove(board Board) (Position, error) {
	emptyCells := getEmptyCells(board)

	if len(emptyCells) == 0 {
		return Position{}, errors.New("no moves left")
	}

	return emptyCells[rand.Intn(len(emptyCells))], nil
}

// RandomMover implements the Mover interface and selects
// a random empty cell from the board for making a move.
type RandomMover struct{}

// Move finds out randoms valid position and returns back
func (r *RandomMover) Move(board Board, _ string) (Position, error) {
	return randomMove(board)
}

// DefensiveMover implements the Mover intereface
type DefensiveMover struct{}

// Move checks for opponent winning move if it didn't find returns back position with random move
func (d *DefensiveMover) Move(board Board, player string) (Position, error) {
	opponent := "O"
	if player == "O" {
		opponent = "X"
	}

	// trying to block opponent
	if pos, ok := findWinningMove(board, opponent); ok {
		return pos, nil
	}

	return randomMove(board)
}

// OffensiveMover implements the Mover interface
type OffensiveMover struct{}

// Move first check for win, if not defensive, if not then random move
func (o *OffensiveMover) Move(board Board, player string) (Position, error) {
	opponent := "O"
	if player == "O" {
		opponent = "X"
	}

	// try to win
	if pos, ok := findWinningMove(board, player); ok {
		return pos, nil
	}

	// try to block
	if pos, ok := findWinningMove(board, opponent); ok {
		return pos, nil
	}

	// random
	return randomMove(board)
}
