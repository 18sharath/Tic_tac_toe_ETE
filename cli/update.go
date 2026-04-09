package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"time"

)

type botMsg struct {
	game *Game
}

func botPlayCmd(id string) tea.Cmd {
	return func() tea.Msg {

		time.Sleep(500 * time.Millisecond) 

		g, err := GetGame(id)
		if err != nil {
			return nil
		}

		return botMsg{game: g}
	}
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case botMsg:

		m.game = msg.game

		// continue loop if game not finished
		if m.game.Winner == "" && !m.game.Draw {
			return m, botPlayCmd(m.game.ID)
		}

		return m, nil

	case tea.KeyMsg:

		switch msg.String() {
		case "r":
			if m.screen == gameScreen {
				g, err := CreateGame(m.mode, m.difficulty)
				if err == nil {

					m.row = 0
					m.col = 0
					m.game = g
				}
				return m, nil
			}

		case "b":
			if m.screen == gameScreen {
				m.screen = menuScreen
				m.cursor = 0
				m.game = nil
				return m, nil
			}

		case "q":
			return m, tea.Quit

		case "up":
			if m.screen == gameScreen && m.row > 0 {
				m.row--
			} else if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			max := len(menuOptions) - 1
			if m.screen == difficultyScreen {
				max = len(difficultyOptions) - 1
			}
			if m.screen == gameScreen && m.row < 2 {
				m.row++
			} else if m.cursor < max {
				m.cursor++
			}

		case "left":
			if m.screen == gameScreen && m.col > 0 {
				m.col--
			}

		case "right":
			if m.screen == gameScreen && m.col < 2 {
				m.col++
			}

		case "enter":

			switch m.screen {

			case menuScreen:
				switch m.cursor{
				case 0:
					m.mode = 1
					g, _ := CreateGame(1, 0)
					m.game = g
					m.screen = gameScreen

				case 1:
					m.mode = 2
					m.screen = difficultyScreen
					m.cursor = 0

				case 2:
					m.mode = 3
					m.screen = difficultyScreen 
					m.cursor = 0

				default :
					return m, tea.Quit
				}

			case difficultyScreen:

				if m.cursor == 3 {
					m.screen = menuScreen
					m.cursor = 0
					return m, nil
				}

				m.difficulty = m.cursor + 1

				g, err := CreateGame(m.mode, m.difficulty)
				if err != nil {
					return m, nil
				}

				m.game = g
				m.row = 0
				m.col = 0

				m.screen = gameScreen 

				if m.mode == 3 {
					return m, botPlayCmd(m.game.ID)
				}

				return m, nil
			
			case gameScreen:

				g, err := MakeMove(m.game.ID, m.game.Turn, m.row, m.col)
				if err == nil {
					m.game = g
				}
			}
		}
	}

	return m, nil
}