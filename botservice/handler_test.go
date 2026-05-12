package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMoveHandlerSuccess(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"O", "O", ""},
			{"X", "X", ""},
			{"", "", ""},
		},
		Player: "O",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp moveResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Row)
	assert.Equal(t, 2, resp.Col)
}

func TestMoveHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/move", nil)
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Contains(t, rec.Body.String(), "method not allowed")
}

func TestMoveHandlerInvalidRequestBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBufferString(`{invalid`))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request body")
}

func TestMoveHandlerEmptyBoard(t *testing.T) {
	body := `{"board": [], "player": "O"}`
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "board must not be empty")
}

func TestMoveHandlerNonSquareBoard(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"X", "O"},
			{"", "", ""},
		},
		Player: "O",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "board must be square (NxN)")
}

func TestMoveHandlerFullBoardNoMovesAvailable(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"X", "O", "X"},
			{"X", "O", "O"},
			{"O", "X", "X"},
		},
		Player: "O",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "no moves available")
}

func TestMoveHandlerDefaultsToOWhenPlayerEmpty(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"O", "O", ""},
			{"X", "X", ""},
			{"", "", ""},
		},
		Player: "",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp moveResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, resp.Row)
}

func TestMoveHandlerAsPlayerX(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"X", "X", ""},
			{"O", "O", ""},
			{"", "", ""},
		},
		Player: "X",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	moveHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp moveResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Row)
	assert.Equal(t, 2, resp.Col)
}

func TestMoveHandlerEncodeError(t *testing.T) {
	body := moveRequest{
		Board: [][]string{
			{"O", "O", ""},
			{"X", "X", ""},
			{"", "", ""},
		},
		Player: "O",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	w := &failingWriter{header: make(http.Header)}

	moveHandler(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
}

func TestValidateBoardSuccess(t *testing.T) {
	board := [][]string{
		{"X", "O", ""},
		{"", "X", ""},
		{"O", "", ""},
	}

	err := validateBoard(board)
	assert.NoError(t, err)
}

func TestValidateBoardEmpty(t *testing.T) {
	err := validateBoard([][]string{})
	assert.EqualError(t, err, "board must not be empty")
}

func TestValidateBoardNonSquare(t *testing.T) {
	board := [][]string{
		{"X", "O"},
		{"", "", ""},
	}

	err := validateBoard(board)
	assert.EqualError(t, err, "board must be square (NxN)")
}

// failingWriter simulates a ResponseWriter whose Write always fails.
type failingWriter struct {
	header     http.Header
	statusCode int
}

func (f *failingWriter) Header() http.Header { return f.header }
func (f *failingWriter) WriteHeader(code int) {
	if f.statusCode == 0 {
		f.statusCode = code
	}
}
func (f *failingWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestNewServer(t *testing.T) {
	srv := newServer(":9091")

	assert.NotNil(t, srv)
	assert.Equal(t, ":9091", srv.Addr)
	assert.NotNil(t, srv.Handler)
	assert.Equal(t, 5*time.Second, srv.ReadTimeout)
	assert.Equal(t, 10*time.Second, srv.WriteTimeout)
	assert.Equal(t, 120*time.Second, srv.IdleTimeout)
}

func TestNewServerRouteMove(t *testing.T) {
	srv := newServer(":9092")

	body := moveRequest{
		Board: [][]string{
			{"O", "O", ""},
			{"X", "X", ""},
			{"", "", ""},
		},
		Player: "O",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
