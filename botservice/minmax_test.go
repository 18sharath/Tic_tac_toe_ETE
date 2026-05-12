package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckWinnerRow(t *testing.T) {
	board := [][]string{
		{"X", "X", "X"},
		{"O", "", "O"},
		{"", "", ""},
	}

	assert.Equal(t, "X", checkWinner(board))
}

func TestCheckWinnerColumn(t *testing.T) {
	board := [][]string{
		{"O", "X", ""},
		{"O", "X", ""},
		{"O", "", ""},
	}

	assert.Equal(t, "O", checkWinner(board))
}

func TestCheckWinnerMainDiagonal(t *testing.T) {
	board := [][]string{
		{"X", "O", ""},
		{"O", "X", ""},
		{"", "", "X"},
	}

	assert.Equal(t, "X", checkWinner(board))
}

func TestCheckWinnerAntiDiagonal(t *testing.T) {
	board := [][]string{
		{"X", "O", "O"},
		{"X", "O", ""},
		{"O", "", "X"},
	}

	assert.Equal(t, "O", checkWinner(board))
}

func TestCheckWinnerNoWinner(t *testing.T) {
	board := [][]string{
		{"X", "O", ""},
		{"", "X", ""},
		{"", "", "O"},
	}

	assert.Equal(t, "", checkWinner(board))
}

func TestIsBoardFullTrue(t *testing.T) {
	board := [][]string{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}

	assert.True(t, isBoardFull(board))
}

func TestIsBoardFullFalse(t *testing.T) {
	board := [][]string{
		{"X", "O", "X"},
		{"X", "", "O"},
		{"O", "X", "X"},
	}

	assert.False(t, isBoardFull(board))
}

func TestBestMoveWinImmediately(t *testing.T) {
	board := [][]string{
		{"O", "O", ""},
		{"X", "X", ""},
		{"", "", ""},
	}

	move := bestMove(board, "O")
	assert.Equal(t, 0, move.Row)
	assert.Equal(t, 2, move.Col)
}

func TestBestMoveBlockOpponent(t *testing.T) {
	board := [][]string{
		{"X", "X", ""},
		{"O", "", ""},
		{"", "", ""},
	}

	move := bestMove(board, "O")
	assert.Equal(t, 0, move.Row)
	assert.Equal(t, 2, move.Col)
}

func TestBestMoveEmptyBoard(t *testing.T) {
	board := [][]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	move := bestMove(board, "O")
	assert.NotEqual(t, -1, move.Row)
	assert.NotEqual(t, -1, move.Col)
}

func TestBestMoveFullBoard(t *testing.T) {
	board := [][]string{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}

	move := bestMove(board, "O")
	assert.Equal(t, -1, move.Row)
	assert.Equal(t, -1, move.Col)
}

func TestMinimaxTerminalWinO(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"O", "O", "O"},
		{"X", "X", ""},
		{"", "", ""},
	}

	score := minimax(board, true, -scoreInf, scoreInf)
	assert.Equal(t, scoreWin, score)
}

func TestMinimaxTerminalWinX(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "X", "X"},
		{"O", "O", ""},
		{"", "", ""},
	}

	score := minimax(board, false, -scoreInf, scoreInf)
	assert.Equal(t, scoreLose, score)
}

func TestMinimaxDraw(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}

	score := minimax(board, true, -scoreInf, scoreInf)
	assert.Equal(t, 0, score)
}

func TestTerminalScoreWinO(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"O", "O", "O"},
		{"X", "X", ""},
		{"", "", ""},
	}

	score, done := terminalScore(board)
	assert.True(t, done)
	assert.Equal(t, scoreWin, score)
}

func TestTerminalScoreWinX(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "X", "X"},
		{"O", "O", ""},
		{"", "", ""},
	}

	score, done := terminalScore(board)
	assert.True(t, done)
	assert.Equal(t, scoreLose, score)
}

func TestTerminalScoreDraw(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "O", "X"},
		{"X", "O", "O"},
		{"O", "X", "X"},
	}

	score, done := terminalScore(board)
	assert.True(t, done)
	assert.Equal(t, 0, score)
}

func TestTerminalScoreNotDone(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "O", ""},
		{"", "", ""},
		{"", "", ""},
	}

	_, done := terminalScore(board)
	assert.False(t, done)
}

func TestMaximizingScore(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"O", "O", ""},
		{"X", "X", ""},
		{"", "", ""},
	}

	score := maximizingScore(board, -scoreInf, scoreInf)
	assert.Equal(t, scoreWin, score)
}

func TestMinimizingScore(t *testing.T) {
	botPlayer = "O"
	board := [][]string{
		{"X", "X", ""},
		{"O", "", ""},
		{"", "", ""},
	}

	score := minimizingScore(board, -scoreInf, scoreInf)
	assert.LessOrEqual(t, score, 0)
}

func TestBestMove4x4Board(t *testing.T) {
	board := [][]string{
		{"O", "O", "O", ""},
		{"X", "X", "X", ""},
		{"", "", "", ""},
		{"", "", "", ""},
	}

	move := bestMove(board, "O")
	assert.Equal(t, 0, move.Row)
	assert.Equal(t, 3, move.Col)
}

func TestOpponent(t *testing.T) {
	assert.Equal(t, "O", opponent("X"))
	assert.Equal(t, "X", opponent("O"))
}

func TestBestMoveAsX(t *testing.T) {
	board := [][]string{
		{"X", "X", ""},
		{"O", "O", ""},
		{"", "", ""},
	}

	move := bestMove(board, "X")
	assert.Equal(t, 0, move.Row)
	assert.Equal(t, 2, move.Col)
}

func TestBestMoveAsXBlocksO(t *testing.T) {
	board := [][]string{
		{"O", "O", ""},
		{"X", "", ""},
		{"", "", ""},
	}

	move := bestMove(board, "X")
	assert.Equal(t, 0, move.Row)
	assert.Equal(t, 2, move.Col)
}
