package game

import (
	"errors"
	"math/rand"
)

type Mover interface {
	Move(board Board,player string) (Position, error)
}

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

func findWinningMove(board Board, symbol string) (Position, bool) {
	emptyCells := getEmptyCells(board)

	for _, pos := range emptyCells {
		if isWinningMove(board, pos.Row, pos.Col, symbol) {
			return pos, true
		}
	}
	return Position{}, false
}

func randomMove(board Board) (Position, error) {
	emptyCells := getEmptyCells(board)

	if len(emptyCells) == 0 {
		return Position{}, errors.New("no moves left")
	}

	return emptyCells[rand.Intn(len(emptyCells))], nil
}


type RandomeMover struct{}

func (r *RandomeMover) Move(board Board,player string) (Position, error) {
	return randomMove(board)
}


type DefensiveMover struct{}

func (d *DefensiveMover) Move(board Board, player string) (Position, error) {

	
	opponent := "O"
	if player=="O"{
		opponent="X"
	}

	// trying to block opponent
	if pos, ok := findWinningMove(board, opponent); ok {
		return pos, nil
	}

	return randomMove(board)
}


type OffensiveMover struct{}

func (o *OffensiveMover) Move(board Board,player string) (Position, error) {


	opponent := "O"
	if player == "O" {
		opponent = "X"
	}

	//try to win
	if pos, ok := findWinningMove(board, player); ok {
		return pos, nil
	}

	//try to block
	if pos, ok := findWinningMove(board, opponent); ok {
		return pos, nil
	}

	//random
	return randomMove(board)
}
