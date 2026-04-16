package main

type screen int

const (
	menuScreen screen = iota
	nameScreen
	sizeScreen
	difficultyScreen
	gameScreen
)

type model struct {
	cursor     int
	screen     screen
	mode       int
	difficultyX int
	difficultyO int
	BoardSize  int
	input      string
	inputMode  string
	player1    string
	player2    string
	game       *Game
	row        int
	col        int
}
