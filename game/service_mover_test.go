package game

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetBotServiceURL(t *testing.T) {
	SetBotServiceURL("http://localhost:9090/move")
	assert.Equal(t, "http://localhost:9090/move", botserviceURL)
}

func TestServiceMoverMoveSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req moveRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "X", req.Player)

		resp := moveResponse{Row: 0, Col: 2}
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(resp)
		require.Nil(t, jsonError)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"X", "X", ""},
		{"O", "O", ""},
		{"", "", ""},
	}

	pos, err := mover.Move(board, "X")

	assert.NoError(t, err)
	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

func TestServiceMoverMoveHTTPError(t *testing.T) {
	SetBotServiceURL("http://invalid-url-that-does-not-exist:99999/move")

	mover := &ServiceMover{}
	board := Board{
		{"X", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	_, err := mover.Move(board, "O")

	assert.Error(t, err)
}

func TestServiceMoverMoveInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{invalid json`))
		require.Nil(t, err)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"X", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	_, err := mover.Move(board, "O")

	assert.Error(t, err)
}

func TestServiceMoverMoveEmptyBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := moveResponse{Row: 1, Col: 1}
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(resp)
		require.Nil(t, jsonError)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	pos, err := mover.Move(board, "X")

	assert.NoError(t, err)
	assert.Equal(t, 1, pos.Row)
	assert.Equal(t, 1, pos.Col)
}

func TestServiceMoverMoveAsO(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req moveRequest
		jsonError := json.NewDecoder(r.Body).Decode(&req)
		require.Nil(t, jsonError)
		assert.Equal(t, "O", req.Player)

		resp := moveResponse{Row: 2, Col: 2}
		w.Header().Set("Content-Type", "application/json")
		jsonError = json.NewEncoder(w).Encode(resp)
		require.Nil(t, jsonError)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"X", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	pos, err := mover.Move(board, "O")

	assert.NoError(t, err)
	assert.Equal(t, 2, pos.Row)
	assert.Equal(t, 2, pos.Col)
}

func TestServiceMoverMoveBoardSent(t *testing.T) {
	var receivedBoard Board

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req moveRequest
		jsonError := json.NewDecoder(r.Body).Decode(&req)
		require.Nil(t, jsonError)
		receivedBoard = req.Board

		resp := moveResponse{Row: 0, Col: 0}
		w.Header().Set("Content-Type", "application/json")
		jsonError = json.NewEncoder(w).Encode(resp)
		require.Nil(t, jsonError)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"X", "O", ""},
		{"", "X", ""},
		{"O", "", ""},
	}

	_, err := mover.Move(board, "X")

	assert.NoError(t, err)
	assert.Equal(t, board, receivedBoard)
}

func TestServiceMoverMoveServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		require.Nil(t, err)
	}))
	defer server.Close()

	SetBotServiceURL(server.URL)

	mover := &ServiceMover{}
	board := Board{
		{"X", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	pos, err := mover.Move(board, "O")

	assert.NoError(t, err)
	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 0, pos.Col)
}
