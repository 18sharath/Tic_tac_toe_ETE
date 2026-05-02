// Package main implements API interaction logic for the Tic Tac Toe CLI.
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// baseURL helps to connects with the backend.
var baseURL string

// CreateGameRequest represents the playload required to create a new game.
type CreateGameRequest struct {
	Mode        int `json:"mode"`
	DifficultyX int `json:"difficultyX"`
	DifficultyO int `json:"difficultyO"`
	BoardSize   int `json:"boardSize"`
}

// Game represents the state of the single game instance.
type Game struct {
	ID     string     `json:"id"`
	Board  [][]string `json:"board"`
	Turn   string     `json:"turn"`
	Winner string     `json:"winner"`
	Draw   bool       `json:"draw"`
}

// CreateGame send the http request to create a new game.
func CreateGame(mode int, diffX, diffO, size int) (g *Game, err error) {
	reqBody := CreateGameRequest{
		Mode:        mode,
		DifficultyX: diffX,
		DifficultyO: diffO,
		BoardSize:   size,
	}


	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	u, err := url.JoinPath(baseURL, "games")
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(u, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("http request failed:", err)
		return nil, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err != nil {
			err = cerr
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return g, nil
}

// GetGame helps to get the game based on its ID.
func GetGame(id string) (g *Game, err error) {
	u, err := url.JoinPath(baseURL, "games", id)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return g, nil
}

// MoveRequest contains the payload required to make a move in a game.
type MoveRequest struct {
	Player string `json:"player"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

// MakeMove helps to make move based on the row and col in a game.
func MakeMove(id string, player string, row, col int) (g *Game, err error) {
	req := MoveRequest{
		Player: player,
		Row:    row,
		Col:    col,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	u, err := url.JoinPath(baseURL, "games", id)

	if err != nil {
		return nil, err
	}

	reqHTTP, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	reqHTTP.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(reqHTTP)

	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	g = &Game{}

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return g, nil
}
