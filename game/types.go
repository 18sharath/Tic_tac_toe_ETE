package game

type Mode int

const (
	ModeHumanVsHuman Mode =iota+1
	ModeHumanVsBot 
	ModeBotVsBot
)

type Difficulty int

const (
	DifficultyEasy Difficulty = iota+1
	DifficultyMedium 
	DifficultyHard
)