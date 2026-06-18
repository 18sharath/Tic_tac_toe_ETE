package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"tic_tac_toe/game"
	"tic_tac_toe/store"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	createGameBody = `{
		"mode": 1,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`
	makeMoveBody = `{
		"player": "X",
		"row": 0,
		"col": 0
	}`
)

type failingResponseWriter struct {
	header     http.Header
	statusCode int
}

func newFailingResponseWriter() *failingResponseWriter {
	return &failingResponseWriter{header: make(http.Header)}
}

func (f *failingResponseWriter) Header() http.Header { return f.header }
func (f *failingResponseWriter) WriteHeader(code int) {
	if f.statusCode == 0 {
		f.statusCode = code
	}
}
func (f *failingResponseWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

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

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(createGameBody))
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
	assert.Contains(t, rec.Body.String(), "invalid game mode")
}

func TestCreateGameHandlerInvalidRequestBody(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(`{"mode":`))
	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request body")
}

func TestCreateGameHandlerLessThanThreeBoardSize(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 1,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 1
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}

	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(3), resp["boardSize"])
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

func TestGetGameHandlerEncodeError(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"encode-id",
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)
	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodGet, "/games/encode-id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "encode-id",
	})

	w := newFailingResponseWriter()
	h.GetGameHandler(w, req)

	assert.Equal(t, http.StatusOK, w.statusCode)
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

	req := httptest.NewRequest(http.MethodPost, "/games/move-id/move", bytes.NewBufferString(makeMoveBody))
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

func TestCreateGameHandlerHumanVsBot(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 2,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(2), resp["mode"])
	assert.Equal(t, float64(3), resp["boardSize"])
	assert.Contains(t, rec.Body.String(), `"turn":"X"`)
}

func TestMakeMoveHandlerHumanVsBot(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"hvb-id",
		3,
		game.ModeHumanVsBot,
		game.DifficultyEasy,
		nil,
		game.NewBotMover(game.DifficultyEasy),
	)

	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodPost, "/games/hvb-id/move", bytes.NewBufferString(makeMoveBody))
	req = mux.SetURLVars(req, map[string]string{
		"id": "hvb-id",
	})

	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "X", resp["turn"])
}

func TestMakeMoveHandlerHumanVsBotNotYourTurn(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame(
		"hvb-turn-id",
		3,
		game.ModeHumanVsBot,
		game.DifficultyEasy,
		nil,
		game.NewBotMover(game.DifficultyEasy),
	)

	_ = h.store.Create(g)

	body := `{
		"player": "O",
		"row": 0,
		"col": 0
	}`

	req := httptest.NewRequest(http.MethodPost, "/games/hvb-turn-id/move", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{
		"id": "hvb-turn-id",
	})

	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "not your turn")
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

func TestCreateGameHandlerBotVsBot(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 3,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(3), resp["mode"])
	assert.Equal(t, float64(3), resp["boardSize"])

	winner, _ := resp["winner"].(string)
	draw, _ := resp["draw"].(bool)
	assert.True(t, winner != "" || draw, "expected game to be finished after BotVsBot creation")
}

func TestCreateGameHandlerBotVsBotDifficulties(t *testing.T) {
	difficulties := []struct {
		diffX int
		diffO int
	}{
		{1, 1},
		{1, 2},
		{2, 3},
		{3, 3},
	}

	for _, d := range difficulties {
		h := setupTestHandler()

		body, _ := json.Marshal(map[string]int{
			"mode":        3,
			"difficultyX": d.diffX,
			"difficultyO": d.diffO,
			"boardSize":   3,
		})

		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		h.CreateGameHandler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		winner, _ := resp["winner"].(string)
		draw, _ := resp["draw"].(bool)
		assert.True(t, winner != "" || draw, "expected finished game for diffX=%d diffO=%d", d.diffX, d.diffO)
	}
}

type mockFailingStore struct {
	store.GameStore
}

func (m *mockFailingStore) Create(_ *game.Game) error {
	return fmt.Errorf("store unavailable")
}

func (m *mockFailingStore) Get(_ string) (*game.Game, bool) { return nil, false }
func (m *mockFailingStore) Delete(_ string) error           { return nil }

// mockFailingDeleteStore is a GameStore whose Delete always returns an error.
type mockFailingDeleteStore struct {
	store.GameStore
}

func (m *mockFailingDeleteStore) Create(_ *game.Game) error       { return nil }
func (m *mockFailingDeleteStore) Get(_ string) (*game.Game, bool) { return nil, false }
func (m *mockFailingDeleteStore) Delete(_ string) error           { return fmt.Errorf("delete failed") }

// mockGetAndFailStore returns a stored game on Get but always errors on Create.
type mockGetAndFailStore struct {
	games map[string]*game.Game
}

func newMockGetAndFailStore() *mockGetAndFailStore {
	return &mockGetAndFailStore{games: make(map[string]*game.Game)}
}

func (m *mockGetAndFailStore) Create(g *game.Game) error {
	if _, exists := m.games[g.ID]; exists {
		return fmt.Errorf("store unavailable")
	}
	m.games[g.ID] = g
	return nil
}

func (m *mockGetAndFailStore) Get(id string) (*game.Game, bool) {
	g, ok := m.games[id]
	return g, ok
}

func (m *mockGetAndFailStore) Delete(_ string) error { return nil }

// mockFailingMover is a Mover whose Move always returns an error.
type mockFailingMover struct{}

func (m *mockFailingMover) Move(_ game.Board, _ string) (game.Position, error) {
	return game.Position{}, fmt.Errorf("mover error")
}

func TestCreateGameHandlerStoreCreateError(t *testing.T) {
	h := NewHandler(&mockFailingStore{})

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(createGameBody))
	rec := httptest.NewRecorder()

	h.CreateGameHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to save game")
}

func TestSaveGameError(t *testing.T) {
	h := NewHandler(&mockFailingStore{})
	g := game.NewGame("save-id", 3, game.ModeHumanVsHuman, game.DifficultyEasy, nil, nil)

	rec := httptest.NewRecorder()
	ok := h.saveGame(rec, g)

	assert.False(t, ok)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to save game")
}

func TestSaveGameSuccess(t *testing.T) {
	h := setupTestHandler()
	g := game.NewGame("save-ok-id", 3, game.ModeHumanVsHuman, game.DifficultyEasy, nil, nil)

	rec := httptest.NewRecorder()
	ok := h.saveGame(rec, g)

	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, rec.Code) // no WriteHeader called by saveGame itself
}

func TestMakeMoveHandlerBotVsBotRejected(t *testing.T) {
	h := setupTestHandler()

	createBody := `{
		"mode": 3,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	h.CreateGameHandler(createRec, createReq)
	assert.Equal(t, http.StatusCreated, createRec.Code)

	var created map[string]interface{}
	err := json.Unmarshal(createRec.Body.Bytes(), &created)
	assert.NoError(t, err)
	gameID := created["id"].(string)

	moveBody := makeMoveBody
	req := httptest.NewRequest(http.MethodPost, "/games/"+gameID+"/move", bytes.NewBufferString(moveBody))
	req = mux.SetURLVars(req, map[string]string{
		"id": gameID,
	})

	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "game already finished")
}

func TestCreateGameHandlerEncodeError(t *testing.T) {
	h := setupTestHandler()

	body := `{
		"mode": 1,
		"difficultyX": 1,
		"difficultyO": 1,
		"boardSize": 3
	}`

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(body))
	w := newFailingResponseWriter()

	h.CreateGameHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.statusCode)
}

func TestGetGameFromRequestFound(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame("req-id", 3, game.ModeHumanVsHuman, game.DifficultyEasy, nil, nil)
	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodGet, "/games/req-id", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "req-id"})

	result, err := h.getGameFromRequest(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "req-id", result.ID)
}

func TestGetGameFromRequestNotFound(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/games/missing", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "missing"})

	result, err := h.getGameFromRequest(req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "game not found")
}

func TestDecodeMoveRequestValid(t *testing.T) {
	body := `{"player":"X","row":1,"col":2}`
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBufferString(body))

	mr, err := decodeMoveRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, "X", mr.Player)
	assert.Equal(t, 1, mr.Row)
	assert.Equal(t, 2, mr.Col)
}

func TestDecodeMoveRequestInvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBufferString(`{invalid`))

	_, err := decodeMoveRequest(req)

	assert.EqualError(t, err, "invalid request")
}

func TestWriteJSONResponseSuccess(t *testing.T) {
	data := map[string]string{"key": "value"}

	rec := httptest.NewRecorder()
	writeJSONResponse(rec, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), `"key":"value"`)
}

func TestWriteJSONResponseEncodeError(t *testing.T) {
	w := newFailingResponseWriter()
	writeJSONResponse(w, http.StatusOK, map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.statusCode)
}

func TestMakeMoveHandlerInvalidRequestBody(t *testing.T) {
	h := setupTestHandler()

	g := game.NewGame("decode-id", 3, game.ModeHumanVsHuman, game.DifficultyEasy, nil, nil)
	_ = h.store.Create(g)

	req := httptest.NewRequest(http.MethodPost, "/games/decode-id/move", bytes.NewBufferString(`{invalid`))
	req = mux.SetURLVars(req, map[string]string{"id": "decode-id"})
	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request")
}

func TestMakeMoveHandlerGameNotFound(t *testing.T) {
	h := setupTestHandler()

	body := `{"player":"X","row":0,"col":0}`
	req := httptest.NewRequest(http.MethodPost, "/games/nonexistent/move", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "game not found")
}

func TestDeleteGameHandlerStoreError(t *testing.T) {
	h := NewHandler(&mockFailingDeleteStore{})

	req := httptest.NewRequest(http.MethodDelete, "/games/some-id", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "some-id"})
	rec := httptest.NewRecorder()

	h.DeleteGameHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to delete game")
}

func TestRunBotGameMoverError(t *testing.T) {
	g := game.NewGame(
		"bot-err-id",
		3,
		game.ModeBotVsBot,
		game.DifficultyEasy,
		&mockFailingMover{},
		&mockFailingMover{},
	)

	runBotGame(g)

	assert.Equal(t, "", g.Winner)
	assert.False(t, g.Draw)
}

func TestMakeMoveHandlerSaveError(t *testing.T) {
	s := newMockGetAndFailStore()
	h := NewHandler(s)

	g := game.NewGame("save-move-id", 3, game.ModeHumanVsHuman, game.DifficultyEasy, nil, nil)
	_ = s.Create(g)

	body := `{"player":"X","row":0,"col":0}`
	req := httptest.NewRequest(http.MethodPost, "/games/save-move-id/move", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"id": "save-move-id"})
	rec := httptest.NewRecorder()

	h.MakeMoveHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to save game")
}
