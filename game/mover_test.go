package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmptyCells(t *testing.T) {
	board := Board{
		{"X", "", "O"},
		{"", "X", ""},
		{"O", "", ""},
	}

	cells := getEmptyCells(board)

	assert.Equal(t, 5, len(cells))
}

func TestIsWinningMoveRow(t *testing.T) {
	board := Board{
		{"X", "X", ""},
		{"O", "", ""},
		{"", "", "O"},
	}

	result := isWinningMove(board, 0, 2, "X")
	assert.True(t, result)
}

func TestIsWinningMoveColumn(t *testing.T) {
	board := Board{
		{"X", "O", ""},
		{"X", "", ""},
		{"", "", "O"},
	}

	result := isWinningMove(board, 2, 0, "X")
	assert.True(t, result)
}

func TestIsWinningMoveDiagonal(t *testing.T) {
	board := Board{
		{"X", "O", ""},
		{"", "X", ""},
		{"", "", ""},
	}

	result := isWinningMove(board, 2, 2, "X")

	assert.True(t, result)
}

func TestFindWinningMove(t *testing.T) {
	board := Board{
		{"O", "O", ""},
		{"X", "", ""},
		{"", "", "X"},
	}

	pos, ok := findWinningMove(board, "O")

	assert.True(t, ok)
	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

func TestRandomMoverMove(t *testing.T) {
	board := Board{
		{"X", "", "O"},
		{"", "", ""},
		{"O", "", ""},
	}

	mover := &RandomMover{}
	pos, err := mover.Move(board, "X")

	assert.NoError(t, err)
	assert.True(t, board[pos.Row][pos.Col] == "")
}


func TestRandomMoverNoMovesLeft(t *testing.T) {
	board := Board{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}

	_, err := randomMove(board)

	assert.Error(t, err)
	assert.Equal(t, "no moves left", err.Error())
}

func TestDefensiveMoverBlocksOpponent(t *testing.T) {
	board := Board{
		{"X", "X", ""},
		{"O", "", ""},
		{"", "", "O"},
	}

	mover := &DefensiveMover{}

	pos, err := mover.Move(board, "O")

	assert.NoError(t, err)
	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

func TestOffensiveMoverWinsFirst(t *testing.T) {
	board := Board{
		{"O", "O", ""},
		{"X", "X", ""},
		{"", "", ""},
	}

	mover := &OffensiveMover{}

	pos, err := mover.Move(board, "O")

	assert.NoError(t, err)

	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

func TestOffensiveMoverBlocksIfNoWin(t *testing.T) {
	board := Board{
		{"X", "X", ""},
		{"O", "", ""},
		{"", "", ""},
	}

	mover := &OffensiveMover{}

	pos, err := mover.Move(board, "O")

	assert.NoError(t, err)

	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

