package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewRouting(t *testing.T) {
	m := model{screen: menuScreen}
	assert.Contains(t, m.View(), "TIC TAC TOE")

	m.screen = nameScreen
	assert.Contains(t, m.View(), "Enter")

	m.screen = sizeScreen
	assert.Contains(t, m.View(), "Board Size")

	m.screen = difficultyScreen
	assert.Contains(t, m.View(), "Difficulty")
}

func TestRenderTurn(t *testing.T) {
	m := model{
		game: &Game{Turn: "X"},
	}

	result := m.renderTurn()
	assert.Contains(t, result, "Turn")
	assert.Contains(t, result, "X")

	m.game.Turn = "O"
	result = m.renderTurn()
	assert.Contains(t, result, "O")
}

func TestRenderGameResultWinner(t *testing.T) {
	m := model{
		game: &Game{Winner: "X"},
	}

	result := m.renderGameResult()

	assert.Contains(t, result, "Winner")
	assert.Contains(t, result, "X")
}

func TestRenderGameResultDraw(t *testing.T) {
	m := model{
		game: &Game{Draw: true},
	}

	result := m.renderGameResult()

	assert.Contains(t, result, "Draw")
}

func TestViewGameScreenTurn(t *testing.T) {
	m := model{
		screen: gameScreen,
		game: &Game{
			Turn: "X",
			Board: [][]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
		},
	}

	result := m.viewGameScreen()

	assert.Contains(t, result, "TIC TAC TOE")
	assert.Contains(t, result, "Turn")
}

func TestViewGameScreenWinner(t *testing.T) {
	m := model{
		screen: gameScreen,
		game: &Game{
			Winner: "O",
			Board: [][]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
		},
	}

	result := m.viewGameScreen()

	assert.Contains(t, result, "Winner")
}

func TestViewMenuScreen(t *testing.T) {
	m := model{
		screen: menuScreen,
		cursor: 0,
	}

	result := m.viewMenuScreen()

	assert.Contains(t, result, "TIC TAC TOE")
	assert.Contains(t, result, "▶")
}

func TestViewDifficultyScreen(t *testing.T) {
	m := model{
		screen: difficultyScreen,
		cursor: 1,
		mode:   int(ModeHumanVsBot),
	}

	result := m.viewDifficultyScreen()

	assert.Contains(t, result, "Difficulty")
	assert.Contains(t, result, "▶")
}

func TestViewSizeScreen(t *testing.T) {
	m := model{
		screen: sizeScreen,
		input:  "4",
	}

	result := m.viewSizeScreen()

	assert.Contains(t, result, "Board Size")
	assert.Contains(t, result, "4")
}

func TestViewNameScreen(t *testing.T) {
	m := model{
		screen:    nameScreen,
		input:     "Sharath",
		inputMode: inputName1,
	}

	result := m.viewNameScreen()

	assert.Contains(t, result, "Enter")
	assert.Contains(t, result, "Sharath")

	m.inputMode = inputName2
	result = m.viewNameScreen()

	assert.Contains(t, result, "Player 2")
}

func TestViewRoutingGameScreen(t *testing.T) {
	m := model{
		screen: gameScreen,
		game: &Game{
			Turn:  "X",
			Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}},
		},
	}

	result := m.View()

	assert.Contains(t, result, "TIC TAC TOE")
}

func TestViewRoutingDefault(t *testing.T) {
	m := model{screen: screen(999)}

	result := m.View()

	assert.Equal(t, "", result)
}

func TestViewDifficultyScreenBotVsBotDiffX(t *testing.T) {
	m := model{
		screen:    difficultyScreen,
		cursor:    0,
		mode:      int(ModeBotVsBot),
		inputMode: inputDiffX,
	}

	result := m.viewDifficultyScreen()

	assert.Contains(t, result, "Bot X")
}

func TestViewDifficultyScreenBotVsBotDiffO(t *testing.T) {
	m := model{
		screen:    difficultyScreen,
		cursor:    0,
		mode:      int(ModeBotVsBot),
		inputMode: inputDiffO,
	}

	result := m.viewDifficultyScreen()

	assert.Contains(t, result, "Bot O")
}

func TestViewBotURLScreen(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        inputMode:   inputBotURLX,
        cursor:      0,
        botServices: []BotService{{Name: "Test Bot", URL: "http://test"}},
    }

    result := m.viewBotURLScreen()

    assert.Contains(t, result, "Bot Service for X")
    assert.Contains(t, result, "Test Bot")
    assert.Contains(t, result, "Custom URL")
}

func TestViewBotURLScreenForO(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        inputMode:   inputBotURLO,
        cursor:      0,
        botServices: []BotService{{Name: "Test Bot", URL: "http://test"}},
    }

    result := m.viewBotURLScreen()

    assert.Contains(t, result, "Bot Service for O")
}

func TestViewBotURLScreenCustomSelected(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        inputMode:   inputBotURLX,
        cursor:      1,
        botServices: []BotService{{Name: "Test Bot", URL: "http://test"}},
        input:       "http://custom",
    }

    result := m.viewBotURLScreen()

    assert.Contains(t, result, "▶ Custom URL")
    assert.Contains(t, result, "http://custom")
}

func TestViewRoutingBotURLScreen(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        inputMode:   inputBotURLX,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    result := m.View()

    assert.Contains(t, result, "Bot Service")
}