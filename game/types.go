package game

// Mode represents the type of game being played.
type Mode int

const (
	// ModeHumanVsHuman represents a game between two human players.
	ModeHumanVsHuman Mode = iota + 1

	// ModeHumanVsBot represents a game between a human and a bot.
	ModeHumanVsBot

	// ModeBotVsBot represents a game between two bots.
	ModeBotVsBot
)

// Difficulty represents the difficulty level of a bot player.
type Difficulty int

const (
	// DifficultyEasy represents the easy bot level.
	DifficultyEasy Difficulty = iota + 1

	// DifficultyMedium represents the medium bot level.
	DifficultyMedium

	// DifficultyHard represents the hard bot level.
	DifficultyHard

	// DifficultyService represents bot moves from external service.
	DifficultyService
)
