package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const baseUrl = "http://localhost:8080"

type CreateGameRequest struct {
	Mode       int `json:"mode"`
	Difficulty int `json:"difficulty"`
	BoardSize  int `json:"boardSize"`
}

type Game struct {
	ID     string     `json:"id"`
	Board  [][]string `json:"board"`
	Turn   string     `json:"turn"`
	Winner string     `json:"winner"`
	Draw   bool       `json:"draw"`
}

func CreateGame(mode int, difficulty int,size int ) (*Game, error) {

	reqBody := CreateGameRequest{
		Mode:       mode,
		Difficulty: difficulty,
		BoardSize:  size,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(baseUrl+"/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var g Game
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}

func GetGame(id string) (*Game, error) {

	resp, err := http.Get(baseUrl + "/games/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var g Game
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}


type MoveRequest struct {
	Player string `json:"player"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

func MakeMove(id string, player string, row, col int) (*Game, error) {

	req := MoveRequest{
		Player: player,
		Row:    row,
		Col:    col,
	}

	data, _ := json.Marshal(req)

	resp, err := http.Post(baseUrl+"/games/"+id, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var g Game
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}

	return &g, nil
}
