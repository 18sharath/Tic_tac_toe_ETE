package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUnreachableURL = "http://127.0.0.1:1"
	testBadURL         = "://bad url"
)

func TestCreateGameSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/games", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		resp := Game{
			ID:   "123",
			Turn: "X",
			Board: [][]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	baseURL = server.URL
	game, err := CreateGame(1, 1, 1, 3)

	assert.NoError(t, err)
	assert.NotNil(t, game)
	assert.Equal(t, "123", game.ID)
	assert.Equal(t, "X", game.Turn)
}

func TestGetGameSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/games/abc", r.URL.Path)

		resp := Game{
			ID:   "abc",
			Turn: "O",
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))

	defer server.Close()

	baseURL = server.URL
	game, err := GetGame("abc")

	assert.NoError(t, err)
	assert.Equal(t, game.ID, "abc")
	assert.Equal(t, game.Turn, "O")
}

func TestMakeMoveSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/games/xyz", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)

		resp := Game{
			ID:   "xyz",
			Turn: "O",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))

	defer server.Close()

	baseURL = server.URL
	game, err := MakeMove("xyz", "X", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, "xyz", game.ID)
	assert.Equal(t, "O", game.Turn)
}

func TestCreateGameHTTPFailure(t *testing.T) {
	baseURL = testUnreachableURL

	game, err := CreateGame(1, 1, 1, 3)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestCreateGameInvalidBaseURL(t *testing.T) {
	baseURL = testBadURL

	game, err := CreateGame(1, 1, 1, 3)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestCreateGameBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json {{{"))
	}))
	defer server.Close()

	baseURL = server.URL
	game, err := CreateGame(1, 1, 1, 3)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestGetGameHTTPFailure(t *testing.T) {
	baseURL = testUnreachableURL

	game, err := GetGame("some-id")

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestGetGameInvalidBaseURL(t *testing.T) {
	baseURL = testBadURL

	game, err := GetGame("some-id")

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestGetGameBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json {{{"))
	}))
	defer server.Close()

	baseURL = server.URL
	game, err := GetGame("some-id")

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestMakeMoveHTTPFailure(t *testing.T) {
	baseURL = testUnreachableURL

	game, err := MakeMove("id", "X", 0, 0)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestMakeMoveInvalidBaseURL(t *testing.T) {
	baseURL = testBadURL

	game, err := MakeMove("id", "X", 0, 0)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestMakeMoveBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json {{{"))
	}))
	defer server.Close()

	baseURL = server.URL
	game, err := MakeMove("id", "X", 0, 0)

	assert.Error(t, err)
	assert.Nil(t, game)
}
