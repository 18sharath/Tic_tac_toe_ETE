package main

import "github.com/charmbracelet/lipgloss"

// View renders the current state of the model as a string.
func (m model) View() string {

	switch m.screen {

	case nameScreen:
		label := "Enter your Name:"

		if m.inputMode == "name2" {
			label = "Enter Player 2 Name:"
		}

		return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center, titleStyle.Render(label+m.input))

	case sizeScreen:
		return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center,
			titleStyle.Render("Enter Board Size (>=3):"+m.input),
		)
	case menuScreen:

		s := "\n"
		s += menuTitleStyle.Render("🎮 TIC TAC TOE") + "\n\n"

		for i, opt := range menuOptions {

			item := menuItemStyle.Render(opt)

			if i == m.cursor {
				item = selectedStyle.Render("▶ " + opt)
			}

			s += item + "\n\n"
		}

		s += helpStyle.Render("↑/↓ to move • Enter to select • q to quit")

		return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center, s)

	case difficultyScreen:

		s := "\n"
		title := "Select Difficulty"

		if m.mode == 3 {
			if m.inputMode == "diffX" {
				title = "Select Difficulty for Bot X"
			} else {
				title = "Select Difficulty for Bot O"
			}
		}
		s += menuTitleStyle.Render(title) + "\n\n"

		for i, opt := range difficultyOptions {

			item := menuItemStyle.Render(opt)

			if i == m.cursor {
				item = selectedStyle.Render("▶ " + opt)
			}

			s += item + "\n\n"
		}

		s += helpStyle.Render("↑/↓ to move • Enter to select • b to back • q to quit")

		return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center, s)

	case gameScreen:

		title := lipgloss.PlaceHorizontal(60, lipgloss.Center, titleStyle.Render("🎮 TIC TAC TOE"))
		s := title + "\n\n"
		s += drawBoard(m.game.Board, m.row, m.col)

		if m.game.Winner != "" || m.game.Draw {

			var msg string

			if m.game.Winner != "" {
				w := m.game.Winner
				if w == "X" {
					w = xStyle.Render("X")
				} else {
					w = oStyle.Render("O")
				}
				msg = "🏆 Winner: " + w
			} else {
				msg = "🤝 It's a Draw!"
			}

			msg += "\n\nPress r to restart • b to menu"

			popup := popupStyle.Render(msg)

			s += "\n\n" + lipgloss.PlaceHorizontal(60, lipgloss.Center, popup)
		} else {
			turn := m.game.Turn
			if turn == "X" {
				turn = xStyle.Render("X")
			} else {
				turn = oStyle.Render("O")
			}
			s += "\nTurn: " + turn
		}

		help := lipgloss.PlaceHorizontal(60, lipgloss.Center, infoStyle.Render("r: restart • b: back • q: quit"))
		s += "\n\n" + help
		width := boardWidth(len(m.game.Board)) + 10
		return lipgloss.Place(width, 20, lipgloss.Center, lipgloss.Center, s)
	}

	return ""
}
