package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard(3)

	assert.Equal(t, 3, len(board))
	assert.Equal(t, 3, len(board[0]))
	assert.Equal(t, "", board[0][0])
	assert.Equal(t, "", board[2][2])
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
