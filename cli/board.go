package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func boardWidth(size int) int {
	cellWidth := 7
	return size * (cellWidth + 1)
}

// drawBoard display the current state of the game.
func drawBoard(board [][]string, row, col int) string {
	s := "\n"

	size := len(board)
	cellWidth := 7

	for i := 0; i < size; i++ {
		s += makePaddingRow(size, cellWidth)
		s += buildRow(board, i, row, col, cellWidth)
		s += makePaddingRow(size, cellWidth)

		if i < size-1 {
			s += makeSeparator(size, cellWidth)
		}
	}

	width := boardWidth(size)
	return lipgloss.Place(width+10, 20, lipgloss.Center, lipgloss.Center, s)
}

// centerCell aligns the given content within a fixed-width cell.
func centerCell(content string, cellWidth int) string {
	contentWidth := 1

	padding := (cellWidth - contentWidth) / 2
	if padding < 0 {
		padding = 0
	}

	left := strings.Repeat(" ", padding)
	right := strings.Repeat(" ", cellWidth-padding-contentWidth)

	return left + content + right
}

// renderCell converts a board value into a styled display string.
func renderCell(value string, i, j, cursorRow, cursorCol int) string {
	display := " "

	switch value {
	case "X":
		display = xStyle.Render("X")
	case "O":
		display = oStyle.Render("O")
	}

	if i == cursorRow && j == cursorCol {
		display = cursorStyle.Render(display)
	}

	return display
}

// buildRow constructs a single row of the board by rendering each cell,
// applying cursor highlighting, and aligning content within fixed-width cells.
func buildRow(board [][]string, i, cursorRow, cursorCol, cellWidth int) string {
	size := len(board[i])
	rowStr := ""

	for j := 0; j < size; j++ {
		cell := renderCell(board[i][j], i, j, cursorRow, cursorCol)
		rowStr += centerCell(cell, cellWidth)

		if j < size-1 {
			rowStr += "|"
		}
	}

	return rowStr + "\n"
}

// makeSeparator creates a horizontal separator line between board rows.
func makeSeparator(size, cellWidth int) string {
	line := ""

	for j := 0; j < size; j++ {
		line += strings.Repeat("-", cellWidth)
		if j < size-1 {
			line += "+"
		}
	}

	return line + "\n"
}

// makePaddingRow generates an empty padding row for the board layout.
func makePaddingRow(size, cellWidth int) string {
	line := ""

	for j := 0; j < size; j++ {
		line += strings.Repeat(" ", cellWidth)
		if j < size-1 {
			line += "|"
		}
	}

	return line + "\n"
}
