package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func startGameServer(t *testing.T) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(Game{ID: "g1", Turn: "X", Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}})
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL
}

func TestUpdateBotMsg(t *testing.T) {
	m := model{game: &Game{ID: "1"}}
	msg := botMsg{game: &Game{ID: "1", Winner: "X"}}

	newModel, cmd := m.Update(msg)

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestUpdateKeyMsg(t *testing.T) {
	m := model{}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}

	newModel, cmd := m.Update(msg)

	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestUpdateUnknownMsg(t *testing.T) {
	m := model{}
	newModel, cmd := m.Update(nil)

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleKeyMsgRestart(t *testing.T) {
	startGameServer(t)
	m := model{screen: gameScreen, mode: int(ModeHumanVsHuman), BoardSize: 3}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	newModel, _ := m.handleKeyMsg(msg)

	assert.NotNil(t, newModel)
}

func TestHandleKeyMsgBack(t *testing.T) {
	m := model{screen: gameScreen, game: &Game{}}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}
	newModel, _ := m.handleKeyMsg(msg)

	updated := newModel.(model)
	assert.Equal(t, menuScreen, updated.screen)
}

func TestHandleKeyMsgUp(t *testing.T) {
	m := model{screen: gameScreen, row: 1}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")}
	newModel, _ := m.handleKeyMsg(msg)

	assert.NotNil(t, newModel)
}

func TestHandleKeyMsgEnterOnMenu(t *testing.T) {
	m := model{screen: menuScreen, cursor: 0}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.handleKeyMsg(msg)

	updated := newModel.(model)
	assert.Equal(t, nameScreen, updated.screen)
}

func TestHandleKeyMsgUnknownKey(t *testing.T) {
	m := model{screen: gameScreen}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("z")}
	newModel, cmd := m.handleKeyMsg(msg)

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleKeyMsgOnNameScreen(t *testing.T) {
	m := model{screen: nameScreen, inputMode: inputName1}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")}
	newModel, _ := m.handleKeyMsg(msg)

	updated := newModel.(model)
	assert.Equal(t, "A", updated.input)
}

func TestHandleInputScreenRunes(t *testing.T) {
	m := model{screen: nameScreen, input: "Al"}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")}
	newModel, _ := m.handleInputScreen(msg)

	updated := newModel.(model)
	assert.Equal(t, "Ali", updated.input)
}

func TestHandleInputScreenBackspace(t *testing.T) {
	m := model{screen: nameScreen, input: "Ab"}

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ := m.handleInputScreen(msg)

	updated := newModel.(model)
	assert.Equal(t, "A", updated.input)
}

func TestHandleInputScreenBackspaceEmpty(t *testing.T) {
	m := model{screen: nameScreen, input: ""}

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ := m.handleInputScreen(msg)

	updated := newModel.(model)
	assert.Equal(t, "", updated.input)
}

func TestHandleInputScreenEnter(t *testing.T) {
	m := model{screen: nameScreen, inputMode: inputName1, input: "X", mode: int(ModeHumanVsHuman)}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.handleInputScreen(msg)

	updated := newModel.(model)
	assert.Equal(t, "X", updated.player1)
}

func TestHandleInputScreenOtherKey(t *testing.T) {
	m := model{screen: nameScreen, input: "A"}

	msg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := m.handleInputScreen(msg)

	updated := newModel.(model)
	assert.Equal(t, "A", updated.input)
}

func TestHandleInputEnterName1HumanVsBot(t *testing.T) {
	m := model{screen: nameScreen, inputMode: inputName1, input: "Alice", mode: int(ModeHumanVsBot)}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, "Alice", updated.player1)
	assert.Equal(t, sizeScreen, updated.screen)
	assert.Equal(t, inputSize, updated.inputMode)
}

func TestHandleInputEnterName2(t *testing.T) {
	m := model{screen: nameScreen, inputMode: inputName2, input: "Bob", mode: int(ModeHumanVsHuman)}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, "Bob", updated.player2)
	assert.Equal(t, sizeScreen, updated.screen)
}

func TestHandleInputEnterSizeInvalid(t *testing.T) {
	m := model{screen: sizeScreen, inputMode: inputSize, input: "abc", mode: int(ModeHumanVsHuman)}

	newModel, cmd := m.handleInputEnter()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleInputEnterSizeTooSmall(t *testing.T) {
	startGameServer(t)
	m := model{screen: sizeScreen, inputMode: inputSize, input: "2", mode: int(ModeHumanVsHuman)}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, 3, updated.BoardSize)
}

func TestStartGameAfterSizeHumanVsHuman(t *testing.T) {
	startGameServer(t)
	m := model{mode: int(ModeHumanVsHuman), BoardSize: 3}

	newModel, _ := m.startGameAfterSize()
	updated := newModel.(model)

	assert.Equal(t, gameScreen, updated.screen)
	assert.NotNil(t, updated.game)
}

func TestStartGameAfterSizeHumanVsBot(t *testing.T) {
	m := model{mode: int(ModeHumanVsBot), BoardSize: 3}

	newModel, _ := m.startGameAfterSize()
	updated := newModel.(model)

	assert.Equal(t, difficultyScreen, updated.screen)
}

func TestStartGameAfterSizeBotVsBot(t *testing.T) {
	m := model{mode: int(ModeBotVsBot), BoardSize: 3}

	newModel, _ := m.startGameAfterSize()
	updated := newModel.(model)

	assert.Equal(t, difficultyScreen, updated.screen)
	assert.Equal(t, inputDiffX, updated.inputMode)
}

func TestStartGameAfterSizeCreateError(t *testing.T) {
	baseURL = testUnreachableURL
	m := model{mode: int(ModeHumanVsHuman), BoardSize: 3}

	newModel, cmd := m.startGameAfterSize()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleMovementDownCursor(t *testing.T) {
	m := model{screen: menuScreen, cursor: 0}

	newModel, _ := m.handleMovement("down")
	updated := newModel.(model)

	assert.Equal(t, 1, updated.cursor)
}

func TestHandleMovementDownDifficultyScreen(t *testing.T) {
	m := model{screen: difficultyScreen, cursor: 0}

	newModel, _ := m.handleMovement("down")
	updated := newModel.(model)

	assert.Equal(t, 1, updated.cursor)
}

func TestHandleMovementDownGameRow(t *testing.T) {
	m := model{
		screen: gameScreen, row: 1,
		game: &Game{Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}},
	}

	newModel, _ := m.handleMovement("down")
	updated := newModel.(model)

	assert.Equal(t, 2, updated.row)
}

func TestHandleMovementLeft(t *testing.T) {
	m := model{screen: gameScreen, col: 1, game: &Game{Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}}}

	newModel, _ := m.handleMovement("left")
	updated := newModel.(model)

	assert.Equal(t, 0, updated.col)
}

func TestHandleMovementRight(t *testing.T) {
	m := model{screen: gameScreen, col: 0, game: &Game{Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}}}

	newModel, _ := m.handleMovement("right")
	updated := newModel.(model)

	assert.Equal(t, 1, updated.col)
}

func TestHandleMovementUpCursorAtZero(t *testing.T) {
	m := model{screen: menuScreen, cursor: 0}

	newModel, _ := m.handleMovement("up")
	updated := newModel.(model)

	assert.Equal(t, 0, updated.cursor)
}

func TestHandleRestartOnGameScreen(t *testing.T) {
	startGameServer(t)
	m := model{screen: gameScreen, mode: int(ModeHumanVsHuman), BoardSize: 3}

	newModel, _ := m.handleRestart()
	updated := newModel.(model)

	assert.NotNil(t, updated.game)
}

func TestHandleRestartNotOnGameScreen(t *testing.T) {
	m := model{screen: menuScreen}

	newModel, _ := m.handleRestart()
	updated := newModel.(model)

	assert.Nil(t, updated.game)
}

func TestHandleEnterOnDifficultyScreen(t *testing.T) {
	startGameServer(t)
	m := model{screen: difficultyScreen, cursor: 0, mode: int(ModeHumanVsBot), BoardSize: 3}

	newModel, _ := m.handleEnter()
	updated := newModel.(model)

	assert.Equal(t, gameScreen, updated.screen)
}

func TestHandleEnterOnGameScreen(t *testing.T) {
	startGameServer(t)
	board := [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}
	m := model{screen: gameScreen, row: 0, col: 0, game: &Game{ID: "g1", Turn: "X", Board: board}}

	newModel, _ := m.handleEnter()

	assert.NotNil(t, newModel)
}

func TestHandleMenuSelectionHumanVsBot(t *testing.T) {
	m := model{screen: menuScreen, cursor: 1}

	newModel, _ := m.handleMenuSelection()
	updated := newModel.(model)

	assert.Equal(t, nameScreen, updated.screen)
	assert.Equal(t, int(ModeHumanVsBot), updated.mode)
	assert.Equal(t, inputName1, updated.inputMode)
}

func TestHandleMenuSelectionQuit(t *testing.T) {
	m := model{screen: menuScreen, cursor: 3}

	_, cmd := m.handleMenuSelection()

	assert.NotNil(t, cmd)
}

func TestHandleDifficultySelectionHumanVsBot(t *testing.T) {
	startGameServer(t)
	m := model{mode: int(ModeHumanVsBot), cursor: 0, BoardSize: 3}

	newModel, _ := m.handleDifficultySelection()
	updated := newModel.(model)

	assert.Equal(t, gameScreen, updated.screen)
	assert.NotNil(t, updated.game)
}

func TestHandleDifficultySelectionHumanVsBotError(t *testing.T) {
	baseURL = testUnreachableURL
	m := model{mode: int(ModeHumanVsBot), cursor: 0, BoardSize: 3}

	newModel, cmd := m.handleDifficultySelection()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleDifficultySelectionBotVsBotDiffX(t *testing.T) {
	m := model{mode: int(ModeBotVsBot), cursor: 1, inputMode: inputDiffX}

	newModel, _ := m.handleDifficultySelection()
	updated := newModel.(model)

	assert.Equal(t, 2, updated.difficultyX)
	assert.Equal(t, inputDiffO, updated.inputMode)
}

func TestHandleDifficultySelectionBotVsBotDiffO(t *testing.T) {
	startGameServer(t)
	m := model{mode: int(ModeBotVsBot), cursor: 0, inputMode: inputDiffO, difficultyX: 1, BoardSize: 3}

	newModel, cmd := m.handleDifficultySelection()
	updated := newModel.(model)

	assert.Equal(t, gameScreen, updated.screen)
	assert.NotNil(t, cmd)
}

func TestHandleDifficultySelectionBotVsBotDiffOError(t *testing.T) {
	baseURL = testUnreachableURL
	m := model{mode: int(ModeBotVsBot), cursor: 0, inputMode: inputDiffO, difficultyX: 1, BoardSize: 3}

	newModel, cmd := m.handleDifficultySelection()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleDifficultySelectionOtherMode(t *testing.T) {
	m := model{mode: int(ModeHumanVsHuman), cursor: 0}

	newModel, cmd := m.handleDifficultySelection()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestBotPlayCmdSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(Game{ID: "bot1", Winner: "X"})
	}))
	defer server.Close()
	baseURL = server.URL

	cmd := botPlayCmd("bot1")
	msg := cmd()

	result, ok := msg.(botMsg)
	assert.True(t, ok)
	assert.Equal(t, "bot1", result.game.ID)
}

func TestBotPlayCmdGetGameError(t *testing.T) {
	baseURL = testUnreachableURL

	cmd := botPlayCmd("any-id")
	msg := cmd()

	assert.Nil(t, msg)
}

func TestHandleInputEnterUnknownScreen(t *testing.T) {
	m := model{screen: difficultyScreen, input: "anything"}

	newModel, cmd := m.handleInputEnter()

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestHandleMenuSelectionHumanVsHuman(t *testing.T) {
	m := model{screen: menuScreen, cursor: 0}

	newModel, _ := m.handleMenuSelection()
	updated := newModel.(model)

	assert.Equal(t, nameScreen, updated.screen)
	assert.Equal(t, int(ModeHumanVsHuman), updated.mode)
	assert.Equal(t, inputName1, updated.inputMode)
}

func TestHandleMenuSelectionBotVsBot(t *testing.T) {
	m := model{screen: menuScreen, cursor: 2}

	newModel, _ := m.handleMenuSelection()
	updated := newModel.(model)

	assert.Equal(t, sizeScreen, updated.screen)
	assert.Equal(t, inputSize, updated.inputMode)
	assert.Equal(t, int(ModeBotVsBot), updated.mode)
}

func TestHandleMovementUp(t *testing.T) {
	m := model{screen: gameScreen, row: 1}

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
	m := model{screen: menuScreen, cursor: 1}

	newModel, _ := m.handleMovement("up")
	updated := newModel.(model)

	assert.Equal(t, 0, updated.cursor)
}

func TestHandleBack(t *testing.T) {
	m := model{screen: gameScreen, cursor: 2, game: &Game{}}

	newModel, _ := m.handleBack()
	updated := newModel.(model)

	assert.Equal(t, menuScreen, updated.screen)
	assert.Equal(t, 0, updated.cursor)
	assert.Nil(t, updated.game)
}

func TestHandleInputEnterName1(t *testing.T) {
	m := model{screen: nameScreen, inputMode: inputName1, input: "Player1", mode: int(ModeHumanVsHuman)}

	newModel, _ := m.handleInputEnter()
	updated := newModel.(model)

	assert.Equal(t, "Player1", updated.player1)
	assert.Equal(t, "", updated.input)
	assert.Equal(t, inputName2, updated.inputMode)
}

func TestHandleInputEnterSize(t *testing.T) {
	m := model{screen: sizeScreen, inputMode: inputSize, input: "5", mode: int(ModeBotVsBot)}

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
	m := model{game: &Game{ID: "1", Winner: "", Draw: false}}

	msg := botMsg{game: &Game{ID: "1", Winner: "", Draw: false}}

	newModel, cmd := m.handleBotMsg(msg)

	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestHandleBotMsgStop(t *testing.T) {
	m := model{}

	msg := botMsg{game: &Game{Winner: "X"}}

	newModel, cmd := m.handleBotMsg(msg)

	assert.NotNil(t, newModel)
	assert.Nil(t, cmd)
}

func TestSetupBotURLScreen(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
        resp := BotServiceConfig{Services: []BotService{{Name: "Bot1", URL: "http://bot1"}}}
        _ = json.NewEncoder(w).Encode(resp)
    }))
    defer server.Close()
    baseURL = server.URL

    m := model{}
    updated := m.setupBotURLScreen(inputBotURLX)

    assert.Equal(t, botURLScreen, updated.screen)
    assert.Equal(t, inputBotURLX, updated.inputMode)
    assert.Equal(t, 0, updated.cursor)
    assert.Len(t, updated.botServices, 1)
}

func TestSetupBotURLScreenWithExistingServices(t *testing.T) {
    m := model{botServices: []BotService{{Name: "Existing", URL: "http://existing"}}}
    updated := m.setupBotURLScreen(inputBotURLO)

    assert.Equal(t, botURLScreen, updated.screen)
    assert.Len(t, updated.botServices, 1)
}

func TestHandleBotURLSelectionFromList(t *testing.T) {
    startGameServer(t)
    m := model{
        mode:        int(ModeHumanVsBot),
        inputMode:   inputBotURLO,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        BoardSize:   3,
        difficultyO: 4,
    }

    newModel, _ := m.handleBotURLSelection()
    updated := newModel.(model)

    assert.Equal(t, "http://bot1", updated.botServiceO)
    assert.Equal(t, gameScreen, updated.screen)
}

func TestHandleBotURLSelectionCustomURL(t *testing.T) {
    startGameServer(t)
    m := model{
        mode:        int(ModeHumanVsBot),
        inputMode:   inputBotURLO,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "http://custom",
        BoardSize:   3,
        difficultyO: 4,
    }

    newModel, _ := m.handleBotURLSelection()
    updated := newModel.(model)

    assert.Equal(t, "http://custom", updated.botServiceO)
}

func TestHandleBotURLSelectionEmptyURL(t *testing.T) {
    m := model{
        inputMode:   inputBotURLO,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "",
    }

    newModel, cmd := m.handleBotURLSelection()

    assert.NotNil(t, newModel)
    assert.Nil(t, cmd)
}

func TestHandleBotURLSelectionForX(t *testing.T) {
    m := model{
        mode:        int(ModeBotVsBot),
        inputMode:   inputBotURLX,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    newModel, _ := m.handleBotURLSelection()
    updated := newModel.(model)

    assert.Equal(t, "http://bot1", updated.botServiceX)
    assert.Equal(t, difficultyScreen, updated.screen)
    assert.Equal(t, inputDiffO, updated.inputMode)
}

func TestHandleBotURLSelectionBotVsBot(t *testing.T) {
    startGameServer(t)
    m := model{
        mode:        int(ModeBotVsBot),
        inputMode:   inputBotURLO,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        BoardSize:   3,
        difficultyX: 4,
        difficultyO: 4,
    }

    newModel, cmd := m.handleBotURLSelection()
    updated := newModel.(model)

    assert.Equal(t, gameScreen, updated.screen)
    assert.NotNil(t, cmd)
}

func TestHandleBotURLScreenRunes(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "http://",
    }

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, "http://test", updated.input)
}

func TestHandleBotURLScreenRunesNotCustom(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "",
    }

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, "", updated.input)
}

func TestHandleBotURLScreenBackspace(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "http://test",
    }

    msg := tea.KeyMsg{Type: tea.KeyBackspace}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, "http://tes", updated.input)
}

func TestHandleBotURLScreenBackspaceEmpty(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        input:       "",
    }

    msg := tea.KeyMsg{Type: tea.KeyBackspace}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, "", updated.input)
}

func TestHandleBotURLScreenEnter(t *testing.T) {
    startGameServer(t)
    m := model{
        screen:      botURLScreen,
        mode:        int(ModeHumanVsBot),
        inputMode:   inputBotURLO,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        BoardSize:   3,
        difficultyO: 4,
    }

    msg := tea.KeyMsg{Type: tea.KeyEnter}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, gameScreen, updated.screen)
}

func TestHandleBotURLScreenUp(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    msg := tea.KeyMsg{Type: tea.KeyUp}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, 0, updated.cursor)
}

func TestHandleBotURLScreenUpAtZero(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    msg := tea.KeyMsg{Type: tea.KeyUp}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, 0, updated.cursor)
}

func TestHandleBotURLScreenDown(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    msg := tea.KeyMsg{Type: tea.KeyDown}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, 1, updated.cursor)
}

func TestHandleBotURLScreenDownAtMax(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    msg := tea.KeyMsg{Type: tea.KeyDown}
    newModel, _ := m.handleBotURLScreen(msg)
    updated := newModel.(model)

    assert.Equal(t, 1, updated.cursor)
}

func TestHandleBotURLScreenBack(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    // Test navigation with unknown key type falls through to default
    msg := tea.KeyMsg{Type: tea.KeyTab}
    newModel, cmd := m.handleBotURLScreen(msg)

    assert.NotNil(t, newModel)
    assert.Nil(t, cmd)
}

func TestHandleKeyMsgBotURLScreen(t *testing.T) {
    m := model{
        screen:      botURLScreen,
        cursor:      1,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
    }

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
    newModel, _ := m.handleKeyMsg(msg)
    updated := newModel.(model)

    assert.Equal(t, "x", updated.input)
}

func TestHandleEnterBotURLScreen(t *testing.T) {
    startGameServer(t)
    m := model{
        screen:      botURLScreen,
        mode:        int(ModeHumanVsBot),
        inputMode:   inputBotURLO,
        cursor:      0,
        botServices: []BotService{{Name: "Bot1", URL: "http://bot1"}},
        BoardSize:   3,
        difficultyO: 4,
    }

    newModel, _ := m.handleEnter()
    updated := newModel.(model)

    assert.Equal(t, gameScreen, updated.screen)
}

func TestHandleDifficultySelectionService(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
        resp := BotServiceConfig{Services: []BotService{{Name: "Bot1", URL: "http://bot1"}}}
        _ = json.NewEncoder(w).Encode(resp)
    }))
    defer server.Close()
    baseURL = server.URL

    m := model{mode: int(ModeHumanVsBot), cursor: 3, BoardSize: 3}

    newModel, _ := m.handleDifficultySelection()
    updated := newModel.(model)

    assert.Equal(t, botURLScreen, updated.screen)
    assert.Equal(t, inputBotURLO, updated.inputMode)
}

func TestHandleDifficultySelectionBotVsBotServiceX(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
        resp := BotServiceConfig{Services: []BotService{{Name: "Bot1", URL: "http://bot1"}}}
        _ = json.NewEncoder(w).Encode(resp)
    }))
    defer server.Close()
    baseURL = server.URL

    m := model{mode: int(ModeBotVsBot), cursor: 3, inputMode: inputDiffX, BoardSize: 3}

    newModel, _ := m.handleDifficultySelection()
    updated := newModel.(model)

    assert.Equal(t, botURLScreen, updated.screen)
    assert.Equal(t, inputBotURLX, updated.inputMode)
}

func TestHandleDifficultySelectionBotVsBotServiceO(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
        resp := BotServiceConfig{Services: []BotService{{Name: "Bot1", URL: "http://bot1"}}}
        _ = json.NewEncoder(w).Encode(resp)
    }))
    defer server.Close()
    baseURL = server.URL

    m := model{mode: int(ModeBotVsBot), cursor: 3, inputMode: inputDiffO, difficultyX: 1, BoardSize: 3}

    newModel, _ := m.handleDifficultySelection()
    updated := newModel.(model)

    assert.Equal(t, botURLScreen, updated.screen)
    assert.Equal(t, inputBotURLO, updated.inputMode)
}