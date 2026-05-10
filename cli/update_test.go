package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func startGameServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Game{ID: "g1", Turn: "X", Board: [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}})
	}))
	t.Cleanup(srv.Close)
	baseURL = srv.URL
	return srv
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
	baseURL = "http://127.0.0.1:1"
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
	baseURL = "http://127.0.0.1:1"
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
	baseURL = "http://127.0.0.1:1"
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	baseURL = "http://127.0.0.1:1"

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
