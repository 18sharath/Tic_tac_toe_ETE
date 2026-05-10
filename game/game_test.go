package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errorMover struct{}

func (e *errorMover) Move(_ Board, _ string) (Position, error) {
	return Position{}, errors.New("mover failure")
}

type occupiedMover struct{}

func (o *occupiedMover) Move(_ Board, _ string) (Position, error) {
	return Position{Row: 0, Col: 0}, nil
}

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

func TestApplyMoveEmptyBoard(t *testing.T) {
	var b Board

	result := ApplyMove(b, Position{Row: 0, Col: 0}, "X")

	assert.Equal(t, b, result)
}

func TestApplyMoveOutOfBounds(t *testing.T) {
	board := NewBoard(3)

	result := ApplyMove(board, Position{Row: 5, Col: 5}, "X")

	// out-of-bounds position should be silently ignored
	for i := range result {
		for j := range result[i] {
			assert.Equal(t, "", result[i][j])
		}
	}
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

func TestPlayTurnGameAlreadyFinishedWinner(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Winner = "X"

	err := g.PlayTurn("X", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "game already finished", err.Error())
}

func TestPlayTurnGameAlreadyFinishedDraw(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Draw = true

	err := g.PlayTurn("X", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "game already finished", err.Error())
}

func TestPlayTurnWrongPlayer(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)

	err := g.PlayTurn("O", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "not your turn", err.Error())
}

func TestPlayTurnHumanVsHumanSuccess(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)

	err := g.PlayTurn("X", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "X", g.Board[0][0])
	assert.Equal(t, "O", g.Turn)
}

func TestPlayTurnHumanMoveInvalidPosition(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)

	err := g.PlayTurn("X", 5, 5)

	assert.Error(t, err)
	assert.Equal(t, "invalid position", err.Error())
}

func TestPlayTurnHumanMoveCellOccupied(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Board[0][0] = "O"

	err := g.PlayTurn("X", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "cell already occupied", err.Error())
}

func TestPlayTurnHumanWinsNoFollowUp(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Board = Board{
		{"X", "X", ""},
		{"O", "O", ""},
		{"", "", ""},
	}

	err := g.PlayTurn("X", 0, 2)

	assert.NoError(t, err)
	assert.Equal(t, "X", g.Winner)
}

func TestPlayTurnDrawNoFollowUp(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Board = Board{
		{"O", "X", "O"},
		{"X", "X", "O"},
		{"X", "O", ""},
	}
	g.Turn = "X"

	err := g.PlayTurn("X", 2, 2)

	assert.NoError(t, err)
	assert.True(t, g.Draw)
	assert.Equal(t, "", g.Winner)
}

func TestPlayTurnHumanVsBotBotMovesAfterHuman(t *testing.T) {
	bot := &RandomMover{}
	g := NewGame("1", 3, ModeHumanVsBot, DifficultyEasy, nil, bot)

	err := g.PlayTurn("X", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "X", g.Board[0][0])

	occupied := 0
	for i := range g.Board {
		for j := range g.Board[i] {
			if g.Board[i][j] != "" {
				occupied++
			}
		}
	}
	assert.Equal(t, 2, occupied)
	assert.Equal(t, "X", g.Turn)
}

func TestPlayTurnBotVsBotFirstTurn(t *testing.T) {
	botX := &RandomMover{}
	botO := &RandomMover{}
	g := NewGame("1", 3, ModeBotVsBot, DifficultyEasy, botX, botO)

	err := g.PlayTurn("X", 0, 0)

	assert.NoError(t, err)
	occupied := 0
	for i := range g.Board {
		for j := range g.Board[i] {
			if g.Board[i][j] != "" {
				occupied++
			}
		}
	}
	assert.Equal(t, 2, occupied)
}

func TestMaketurnNoMover(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)

	err := g.Maketurn()

	assert.Error(t, err)
	assert.Equal(t, "no mover available for current player", err.Error())
}

func TestMaketurnMoverReturnsError(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsBot, DifficultyEasy, &errorMover{}, nil)

	err := g.Maketurn()

	assert.Error(t, err)
	assert.Equal(t, "mover failure", err.Error())
}

func TestMaketurnMoverReturnsOccupiedCell(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsBot, DifficultyEasy, &occupiedMover{}, nil)
	g.Board[0][0] = "O"

	err := g.Maketurn()

	assert.Error(t, err)
	assert.Equal(t, "invalid move", err.Error())
}

func TestPlayTurnBotFirstMoveError(t *testing.T) {
	g := NewGame("1", 3, ModeBotVsBot, DifficultyEasy, &errorMover{}, &RandomMover{})

	err := g.PlayTurn("X", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "mover failure", err.Error())
}

func TestPlayTurnFollowUpBotError(t *testing.T) {
	g := NewGame("1", 3, ModeHumanVsBot, DifficultyEasy, nil, &errorMover{})

	err := g.PlayTurn("X", 0, 0)

	assert.Error(t, err)
	assert.Equal(t, "mover failure", err.Error())
}

func TestCloneNilGame(t *testing.T) {
	var g *Game
	result := g.Clone()

	assert.Nil(t, result)
}

func TestCloneDeepCopiesBoard(t *testing.T) {
	g := NewGame("clone-1", 3, ModeHumanVsHuman, DifficultyEasy, nil, nil)
	g.Board[0][0] = "X"

	cloned := g.Clone()

	assert.NotNil(t, cloned)
	assert.Equal(t, "X", cloned.Board[0][0])

	// mutating clone must not affect original
	cloned.Board[0][0] = "O"
	assert.Equal(t, "X", g.Board[0][0])
}

func TestCloneCopiesFields(t *testing.T) {
	bot := &RandomMover{}
	g := NewGame("clone-2", 3, ModeHumanVsBot, DifficultyMedium, nil, bot)
	g.Turn = "O"
	g.Winner = "O"
	g.Draw = true

	cloned := g.Clone()

	assert.Equal(t, g.ID, cloned.ID)
	assert.Equal(t, g.Turn, cloned.Turn)
	assert.Equal(t, g.Winner, cloned.Winner)
	assert.True(t, cloned.Draw)
	assert.Equal(t, g.Mode, cloned.Mode)
	assert.Equal(t, g.Difficulty, cloned.Difficulty)
	assert.Equal(t, g.PlayerO, cloned.PlayerO)
}

func TestCloneNilBoard(t *testing.T) {
	g := &Game{
		ID:    "nil-board",
		Board: nil,
	}

	cloned := g.Clone()

	assert.NotNil(t, cloned)
	assert.Nil(t, cloned.Board)
}
