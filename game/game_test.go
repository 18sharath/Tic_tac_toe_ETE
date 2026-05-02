package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard(3)

	assert.Equal(t, 3, len(board))
	assert.Equal(t, 3, len(board[0]))
	assert.Equal(t, "", board[0][0])
	assert.Equal(t, "", board[2][2])
}

func TestApplyMoveReturnsNewBoard(t *testing.T) {
	board := NewBoard(3)
	board[0][0] = "X"

	updated := ApplyMove(board, Position{Row: 1, Col: 1}, "O")

	assert.Equal(t, "X", board[0][0])
	assert.Equal(t, "", board[1][1])
	assert.Equal(t, "O", updated[1][1])

	updated[0][0] = "Z"
	assert.Equal(t, "X", board[0][0])
}

func TestNewGame(t *testing.T) {
	g := NewGame(
		"game-1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	assert.Equal(t, "game-1", g.ID)
	assert.Equal(t, "X", g.Turn)
	assert.Equal(t, "", g.Winner)
	assert.False(t, g.Draw)
	assert.Equal(t, 3, len(g.Board))
}

func TestMakeMoveSuccess(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	err := g.MakeMove("X", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "X", g.Board[0][0])
	assert.Equal(t, "O", g.Turn)
}

func TestMakeMoveDoesNotMutateOriginalBoardAlias(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	original := g.Board

	err := g.MakeMove("X", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "", original[0][0])
	assert.Equal(t, "X", g.Board[0][0])

	g.Board[0][0] = "O"
	assert.Equal(t, "", original[0][0])
}

func TestMakeMoveWrongTurn(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	err := g.MakeMove("O", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "not your turn", err.Error())
}

func TestEvaluateRowWinner(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	g.Board = Board{
		{"X", "X", "X"},
		{"", "O", ""},
		{"O", "", ""},
	}

	g.Evaluate()

	assert.Equal(t, "X", g.Winner)
	assert.False(t, g.Draw)
}

func TestEvaluateColumnWinner(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	g.Board = Board{
		{"O", "X", "X"},
		{"O", "X", ""},
		{"O", "", ""},
	}

	g.Evaluate()

	assert.Equal(t, "O", g.Winner)
	assert.False(t, g.Draw)
}

func TestEvaluateDiagonalWinner(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	g.Board = Board{
		{"O", "X", "X"},
		{"", "O", ""},
		{"X", "", "O"},
	}

	g.Evaluate()

	assert.Equal(t, "O", g.Winner)
	assert.False(t, g.Draw)
}

func TestEvaluateAntiDiagonalWinner(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	g.Board = Board{
		{"", "O", "X"},
		{"O", "X", ""},
		{"X", "", "O"},
	}

	g.Evaluate()

	assert.Equal(t, "X", g.Winner)
	assert.False(t, g.Draw)
}

func TestEvaluateDraw(t *testing.T) {
	g := NewGame(
		"1",
		3,
		ModeHumanVsHuman,
		DifficultyEasy,
		nil,
		nil,
	)

	g.Board = Board{
		{"X", "O", "X"},
		{"O", "O", "X"},
		{"X", "X", "O"},
	}

	g.Evaluate()
	assert.True(t, g.Draw)
	assert.Equal(t, "", g.Winner)
}

func TestNewBotMover(t *testing.T) {
	m1 := NewBotMover(DifficultyEasy)
	m2 := NewBotMover(DifficultyMedium)
	m3 := NewBotMover(DifficultyHard)

	assert.IsType(t, &RandomMover{}, m1)
	assert.IsType(t, &DefensiveMover{}, m2)
	assert.IsType(t, &OffensiveMover{}, m3)
}
