package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"tic_tac_toe/game"
	"tic_tac_toe/store"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler() *Handler {
	s := store.NewMemoryStore()
	return NewHandler(s)
}

func TestValidatePlayer(t *testing.T) {
	err := validatePlayer("Z")
	assert.Error(t, err)
	assert.Equal(t, "invalid player", err.Error())

	err = validatePlayer("X")
	assert.NoError(t, err)

	err = validatePlayer("O")
	assert.NoError(t, err)
}

func TestCreateGameHandlerSuccess(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 1,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	assert.Contains(t, rec.Body.String(), `"turn":"X"`)
}

func TestCreateGameHandlerInvalidMethod(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestCreateGameHandlerInvalidMode(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 99,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid mode")
}

func TestGetGameHandlerSuccess(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"test-id",
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)

	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodGet, "/games/test-d", nil)

	req = mux.SetURLVars(req, map[string]string{
		"id": "test-id",
	})
	rec := httptest.NewRecorder()

	h.GetGameHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":"test-id"`)
}

func TestGetGameHandlerNotFound(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/games/not-found", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "not-found",
	})

	rec := httptest.NewRecorder()

	h.GetGameHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMakeMoveHandlerSuccess(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"move-id",
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)

	_ = h.store.Create(g)

	body := `{
		"player": "X",
		"row": 0,
		"col": 0
	}`

	req := httptest.NewRequest(http.MethodPost, "/games/move-id/move", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{
		"id": "move-id",
	})

	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"turn":"O"`)
}

func TestMakeMoveHandlerInvalidPlayer(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"move-id",
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)

	_ = h.store.Create(g)

	body := `{
		"player": "Z",
		"row": 0,
		"col": 0
	}`

	req := httptest.NewRequest(http.MethodPost, "/games/move-id/move", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{
		"id": "move-id",
	})

	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid player")
}

func TestDeleteGameHandler(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"delete-id",
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)

	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodDelete, "/games/delete-id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "delete-id",
	})

	rec := httptest.NewRecorder()

	h.DeleteGameHandler(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
