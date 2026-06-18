package main

type screen int

const (
	menuScreen screen = iota
	nameScreen
	sizeScreen
	difficultyScreen
	botURLScreen
	gameScreen
)

const (
	inputName1 = "name1"
	inputName2 = "name2"
	inputSize  = "size"
	inputDiffX = "diffX"
	inputDiffO = "diffO"
	inputBotURLX = "botURLX"
	inputBotURLO = "botURLO"
)

type mode int

const (
	ModeHumanVsHuman mode = iota + 1
	ModeHumanVsBot
	ModeBotVsBot
)

type BotService struct{
	Name string `json:"name"`
	URL string 	`json:"url"`
	Description string `json:"description"`
}

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
	botServiceX string
	botServiceO string
	botServices []BotService
}
