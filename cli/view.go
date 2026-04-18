package main

import "github.com/charmbracelet/lipgloss"



func (m model) View() string {
	switch m.screen {
	case nameScreen:
		return m.viewNameScreen()
	case sizeScreen:
		return m.viewSizeScreen()
	case menuScreen:
		return m.viewMenuScreen()
	case difficultyScreen:
		return m.viewDifficultyScreen()
	case gameScreen:
		return m.viewGameScreen()
	}
	return ""
}


// renderTurn shows the current player's turn.
func (m model) renderTurn() string {
	turn := m.game.Turn

	if turn == "X" {
		turn = xStyle.Render("X")
	} else {
		turn = oStyle.Render("O")
	}

	return "\nTurn: " + turn
}

// renderGameResult displays the winner or draw message.
func (m model) renderGameResult() string {
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

	return "\n\n" + lipgloss.PlaceHorizontal(60, lipgloss.Center, popup)
}

// viewGameScreen renders the active game board along with status and controls.
func (m model) viewGameScreen() string {
	title := lipgloss.PlaceHorizontal(60, lipgloss.Center, titleStyle.Render("🎮 TIC TAC TOE"))
	s := title + "\n\n"

	s += drawBoard(m.game.Board, m.row, m.col)

	if m.game.Winner != "" || m.game.Draw {
		s += m.renderGameResult()
	} else {
		s += m.renderTurn()
	}

	help := lipgloss.PlaceHorizontal(60, lipgloss.Center,
		infoStyle.Render("r: restart • b: back • q: quit"),
	)

	s += "\n\n" + help

	width := boardWidth(len(m.game.Board)) + 10
	return lipgloss.Place(width, 20, lipgloss.Center, lipgloss.Center, s)
}

// viewDifficultyScreen renders difficulty selection UI for bots.
func (m model) viewDifficultyScreen() string {
	s := "\n"
	title := "Select Difficulty"

	if m.mode == int(ModeBotVsBot) {
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
}


// viewMenuScreen renders the main menu with selectable options.
func (m model) viewMenuScreen() string {
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
}

// viewSizeScreen renders input UI for board size selection.
func (m model) viewSizeScreen() string {
	return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center,
		titleStyle.Render("Enter Board Size (>=3):"+m.input),
	)
}

// viewNameScreen renders input UI for player names.
func (m model) viewNameScreen() string {
	label := "Enter your Name:"

	if m.inputMode == "name2" {
		label = "Enter Player 2 Name:"
	}

	return lipgloss.Place(60, 20, lipgloss.Center, lipgloss.Center,
		titleStyle.Render(label+m.input),
	)
}