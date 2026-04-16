package game

// Mode represents the type of game being played.
type Mode int

const (
	ModeHumanVsHuman Mode =iota+1
	ModeHumanVsBot 
	ModeBotVsBot
)

// Difficulty represents the difficulty level of a bot player.
type Difficulty int

const (
	DifficultyEasy Difficulty = iota+1
	DifficultyMedium 
	DifficultyHard
)