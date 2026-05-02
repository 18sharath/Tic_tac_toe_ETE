package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestHandleMenuSelectionHumanVsHuman(t *testing.T) {
	m := model{
		screen: menuScreen,
		cursor: 0,
	}

	newModel, _ := m.handleMenuSelection()

	updated := newModel.(model)

	assert.Equal(t, nameScreen, updated.screen)
	assert.Equal(t, int(ModeHumanVsHuman), updated.mode)
	assert.Equal(t, inputName1, updated.inputMode)
}

func TestHandleMenuSelectionBotVsBot(t *testing.T) {
	m := model{
		screen: menuScreen,
		cursor: 2,
	}

	newModel, _ := m.handleMenuSelection()
	updated := newModel.(model)

	assert.Equal(t, sizeScreen, updated.screen)
	assert.Equal(t, inputSize, updated.inputMode)
	assert.Equal(t, int(ModeBotVsBot), updated.mode)
}

func TestHandleMovementUp(t *testing.T) {
	m := model{
		screen: gameScreen,
		row:    1,
	}

	newModel, _ := m.handleMovement("up")
	updated := newModel.(model)

	assert.Equal(t, 0, updated.row)
}

func TestHandleMovementDown(t *testing.T) {
	m := model{
		screen: gameScreen,
		row:    0,
		game: &Game{
			Board: [][]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
		},
	}

	newModel, _ := m.handleMovement("down")
	updated := newModel.(model)

	assert.Equal(t, 1, updated.row)
}

func TestHandleMovementCursor(t *testing.T) {
	m := model{
		screen: menuScreen,
		cursor: 1,
	}

	newModel, _ := m.handleMovement("up")
	updated := newModel.(model)

	assert.Equal(t, 0, updated.cursor)
}

func TestHandleBack(t *testing.T) {
	m := model{
		screen: gameScreen,
		cursor: 2,
		game:   &Game{},
	}

	newModel, _ := m.handleBack()
	updated := newModel.(model)

	assert.Equal(t, menuScreen, updated.screen)
	assert.Equal(t, 0, updated.cursor)
	assert.Nil(t, updated.game)
}

func TestHandleInputEnterName1(t *testing.T) {
	m := model{
		screen:    nameScreen,
		inputMode: inputName1,
		input:     "Player1",
		mode:      int(ModeHumanVsHuman),
	}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, "Player1", updated.player1)
	assert.Equal(t, "", updated.input)
	assert.Equal(t, inputName2, updated.inputMode)
}

func TestHandleInputEnterSize(t *testing.T) {
	m := model{
		screen:    sizeScreen,
		inputMode: inputSize,
		input:     "5",
		mode:      int(ModeBotVsBot),
	}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, 5, updated.BoardSize)
	assert.Equal(t, difficultyScreen, updated.screen)
}

func TestHandleKeyMsgQuit(t *testing.T) {
	m := model{}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}

	newModel, cmd := m.handleKeyMsg(msg)

	assert.NotNil(t, cmd)
	assert.NotNil(t, newModel)
}

func TestHandleBotMsgContinue(t *testing.T) {
	m := model{
		game: &Game{
			ID:     "1",
			Winner: "",
			Draw:   false,
		},
	}

	msg := botMsg{
		game: &Game{
			ID:     "1",
			Winner: "",
			Draw:   false,
		},
	}

	newModel, cmd := m.handleBotMsg(msg)

	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestHandleBotMsgStop(t *testing.T) {
	m := model{}

	msg := botMsg{
		game: &Game{
			Winner: "X",
		},
	}

	newModel, cmd := m.handleBotMsg(msg)

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}
