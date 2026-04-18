package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
// Update processes a Bubble Tea message and returns the updated model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case botMsg:
		return m.handleBotMsg(msg)

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m model) handleBotMsg(msg botMsg) (tea.Model, tea.Cmd) {
	m.game = msg.game

	if m.game.Winner == "" && !m.game.Draw {
		return m, botPlayCmd(m.game.ID)
	}

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// input screens handled separately
	if m.screen == nameScreen || m.screen == sizeScreen {
		return m.handleInputScreen(msg)
	}

	switch msg.String() {
	case "r":
		return m.handleRestart()
	case "b":
		return m.handleBack()
	case "q":
		return m, tea.Quit
	case "up", "down", "left", "right":
		return m.handleMovement(msg.String())
	case "enter":
		return m.handleEnter()
	}

	return m, nil
}

func (m model) handleInputScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.Type {
	case tea.KeyRunes:
		m.input += string(msg.Runes)

	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	case tea.KeyEnter:
		return m.handleInputEnter()
	}

	return m, nil
}


func (m model) handleInputEnter() (tea.Model, tea.Cmd) {
	// NAME FLOW
	if m.screen == nameScreen {

		if m.inputMode == inputName1 {
			m.player1 = m.input
			m.input = ""

			if m.mode == 1 {
				m.inputMode = inputName2
			} else {
				m.screen = sizeScreen
				m.inputMode = inputSize
			}
			return m, nil
		}

		if m.inputMode == inputName2 {
			m.player2 = m.input
			m.input = ""
			m.screen = sizeScreen
			m.inputMode = inputSize
			return m, nil
		}
	}

	// SIZE FLOW
	if m.screen == sizeScreen {
		size := 3

		n, err := fmt.Sscanf(m.input, "%d", &size)
		if err != nil || n != 1 {
			return m, nil
		}

		if size < 3 {
			size = 3
		}

		m.BoardSize = size
		m.input = ""

		return m.startGameAfterSize()
	}

	return m, nil
}


func (m model) startGameAfterSize() (tea.Model, tea.Cmd) {

	if m.mode == int(ModeHumanVsBot) {
		m.screen = difficultyScreen
		m.cursor = 0
		return m, nil
	}

	if m.mode == int(ModeBotVsBot) {
		m.screen = difficultyScreen
		m.cursor = 0
		m.inputMode = inputDiffX
		return m, nil
	}

	g, err := CreateGame(m.mode, m.difficultyX, m.difficultyO, m.BoardSize)
	if err != nil {
		return m, nil
	}

	m.game = g
	m.screen = gameScreen
	m.row, m.col = 0, 0

	return m, nil
}

func (m model) handleMovement(key string) (tea.Model, tea.Cmd) {

	switch key {

	case "up":
		if m.screen == gameScreen && m.row > 0 {
			m.row--
		} else if m.cursor > 0 {
			m.cursor--
		}

	case "down":
		maxCursor := len(menuOptions) - 1
		if m.screen == difficultyScreen {
			maxCursor = len(difficultyOptions) - 1
		}

		if m.screen == gameScreen && m.row < len(m.game.Board)-1 {
			m.row++
		} else if m.cursor < maxCursor {
			m.cursor++
		}

	case "left":
		if m.screen == gameScreen && m.col > 0 {
			m.col--
		}

	case "right":
		if m.screen == gameScreen && m.col < len(m.game.Board)-1 {
			m.col++
		}
	}

	return m, nil
}

func (m model) handleRestart() (tea.Model, tea.Cmd) {
	if m.screen == gameScreen {
		g, err := CreateGame(m.mode, m.difficultyX, m.difficultyO, m.BoardSize)
		if err == nil {
			m.game = g
			m.row, m.col = 0, 0
		}
	}
	return m, nil
}

func (m model) handleBack() (tea.Model, tea.Cmd) {
	if m.screen == gameScreen || m.screen == difficultyScreen {
		m.screen = menuScreen
		m.cursor = 0
		m.game = nil
	}
	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {

	switch m.screen {

	case menuScreen:
		return m.handleMenuSelection()

	case difficultyScreen:
		return m.handleDifficultySelection()

	case gameScreen:
		g, err := MakeMove(m.game.ID, m.game.Turn, m.row, m.col)
		if err == nil {
			m.game = g
		}
	}

	return m, nil
}

func (m model) handleMenuSelection() (tea.Model, tea.Cmd) {

	switch m.cursor {
	case 0:
		m.mode = int(ModeHumanVsHuman)
	case 1:
		m.mode = int(ModeHumanVsBot)
	case 2:
		m.mode = int(ModeBotVsBot)
	default:
		return m, tea.Quit
	}

	m.screen = nameScreen
	m.inputMode = inputName1
	m.input = ""

	if m.mode == int(ModeBotVsBot) {
		m.screen = sizeScreen
		m.inputMode = inputSize
	}

	return m, nil
}

func (m model) handleDifficultySelection() (tea.Model, tea.Cmd) {

	diff := m.cursor + 1

	if m.mode == int(ModeHumanVsBot) {
		m.difficultyO = diff

		g, err := CreateGame(m.mode, 0, m.difficultyO, m.BoardSize)
		if err != nil {
			return m, nil
		}

		m.game = g
		m.screen = gameScreen
		m.row, m.col = 0, 0

		return m, nil
	}

	if m.mode == int(ModeBotVsBot) {

		if m.inputMode == inputDiffX {
			m.difficultyX = diff
			m.inputMode = inputDiffO
			return m, nil
		}

		if m.inputMode == inputDiffO {
			m.difficultyO = diff

			g, err := CreateGame(m.mode, m.difficultyX, m.difficultyO, m.BoardSize)
			if err != nil {
				return m, nil
			}

			m.game = g
			m.screen = gameScreen
			m.row, m.col = 0, 0

			return m, botPlayCmd(m.game.ID)
		}
	}

	return m, nil
}


