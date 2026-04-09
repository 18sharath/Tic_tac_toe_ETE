package main

type screen int

const (
	menuScreen screen = iota
	difficultyScreen
	gameScreen
)

type model struct {
	cursor     int
	screen     screen
	mode       int
	difficulty int

	game   *Game
	row    int
	col    int
	status string
}
