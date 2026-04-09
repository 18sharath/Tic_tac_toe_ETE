package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func drawBoard(board [][]string, row, col int) string {
	s := "\n"

	for i := 0; i < 3; i++ {

		s += "       |       |       \n"

		for j := 0; j < 3; j++ {

			cell := board[i][j]

			switch cell{
			case "":
				cell = " "

			case "X":
				cell = xStyle.Render("X")

			case "O":
				cell = oStyle.Render("O")
			}

			// cursor highlight
			if i == row && j == col {
				cell = cursorStyle.Render(cell)
			}

			s += fmt.Sprintf("   %s   ", cell)

			if j < 2 {
				s += "|"
			}
		}

		s += "\n"

		// bottom padding line
		s += "       |       |       \n"

		if i < 2 {
			s += "-------+-------+-------\n"
		}
	}

	return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center, s)
}
