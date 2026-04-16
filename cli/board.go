package main

import (
	"strings"
	"github.com/charmbracelet/lipgloss"
)

func boardWidth(size int) int {
	cellWidth :=7
	return size*(cellWidth+1)
}

// drawBoard display the current state of the game.
func drawBoard(board [][]string, row, col int) string {
	s := "\n"

	size := len(board)

	cellWidth := 7 

	makePaddingRow := func() string {
		line := ""
		for j := 0; j < size; j++ {
			line += strings.Repeat(" ", cellWidth)
			if j < size-1 {
				line += "|"
			}
		}
		return line + "\n"
	}

	makeSeparator := func() string {
		line := ""
		for j := 0; j < size; j++ {
			line += strings.Repeat("-", cellWidth)
			if j < size-1 {
				line += "+"
			}
		}
		return line + "\n"
	}

	for i := 0; i < size; i++ {

		s += makePaddingRow()

		for j := 0; j < size; j++ {
			raw := board[i][j]
			display := " "

			switch raw{
			case "X":
				display = xStyle.Render("X")
			case "O":
				display = oStyle.Render("O")
			}

			if i == row && j == col {
				display = cursorStyle.Render(display)
			}

			
			contentWidth := 1 
			padding := (cellWidth - contentWidth) / 2

			if padding < 0 {
				padding = 0 
			}

			left := strings.Repeat(" ", padding)
			right := strings.Repeat(" ", cellWidth-padding-contentWidth)

			s += left + display + right

			if j < size-1 {
				s += "|"
			}
		}

		s += "\n"

		s += makePaddingRow()

		if i < size-1 {
			s += makeSeparator()
		}
	}
	width := boardWidth(len(board))
	return lipgloss.Place(width+10, 20, lipgloss.Center, lipgloss.Center, s)
}
