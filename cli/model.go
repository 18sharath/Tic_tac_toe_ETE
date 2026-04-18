package main

type screen int

const (
	menuScreen screen = iota
	nameScreen
	sizeScreen
	difficultyScreen
	gameScreen
)

const (
	inputName1 = "name1"
	inputName2 = "name2"
	inputSize  = "size"
	inputDiffX = "diffX"
	inputDiffO = "diffO"
)

type mode int

const (
	ModeHumanVsHuman mode = iota+1
	ModeHumanVsBot  
	ModeBotVsBot     
)

type model struct {
	cursor      int
	screen      screen
	mode        int
	difficultyX int
	difficultyO int
	BoardSize   int
	input       string
	inputMode   string
	player1     string
	player2     string
	game        *Game
	row         int
	col         int
}
