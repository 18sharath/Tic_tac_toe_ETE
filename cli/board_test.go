package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardWidth(t *testing.T) {
	width := boardWidth(3)
	assert.Equal(t, 3*(7+1), width)
}

func TestCenterCell(t *testing.T) {
	result := centerCell("X", 7)

	assert.Contains(t, result, "X")
	assert.Equal(t, 7, len(result))
}

func TestRenderCellEmpty(t *testing.T) {
	result := renderCell("", 0, 0, 1, 1)

	assert.Equal(t, " ", result)
}

func TestRenderCellWithX(t *testing.T) {
	result := renderCell("X", 0, 0, 1, 1)

	assert.Contains(t, result, "X")
}

func TestRenderCellWithCursor(t *testing.T) {
	result := renderCell("X", 0, 0, 0, 0)

	assert.Contains(t, result, "X")
}

func TestBuildRow(t *testing.T) {
	board := [][]string{
		{"X", "O", ""},
	}

	row := buildRow(board, 0, -1, -1, 7)

	assert.Contains(t, row, "X")
	assert.Contains(t, row, "O")
	assert.Contains(t, row, "|")
}

func TestMakeSeparator(t *testing.T) {
	sep := makeSeparator(3, 7)

	assert.Contains(t, sep, "+")
	assert.Contains(t, sep, "-")
}

func TestMakePaddingRow(t *testing.T) {
	row := makePaddingRow(3, 7)

	assert.Contains(t, row, "|")
}

func TestDrawBoardBasic(t *testing.T) {
	board := [][]string{
		{"X", "O", ""},
		{"", "X", ""},
		{"O", "", ""},
	}

	result := drawBoard(board, 0, 0)

	assert.Contains(t, result, "X")
	assert.Contains(t, result, "O")
	assert.Contains(t, result, "|")
}

func TestDrawBoardDynamicSize(t *testing.T) {
	board := [][]string{
		{"X", "O", "", ""},
		{"", "X", "", ""},
		{"O", "", "X", ""},
		{"", "", "", ""},
	}

	result := drawBoard(board, 0, 0)

	assert.Contains(t, result, "X")
	assert.Contains(t, result, "O")
	assert.True(t, strings.Contains(result, "+"))
}