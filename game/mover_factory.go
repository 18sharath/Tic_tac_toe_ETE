package game

// NewBotMover returns a Mover implementation based on the given difficulty level.
func NewBotMover(difficulty Difficulty) Mover {
	switch difficulty {
	case DifficultyEasy:
		return &RandomMover{}
	case DifficultyMedium:
		return &DefensiveMover{}
	case DifficultyHard:
		return &OffensiveMover{}
	case DifficultyService:
		return &ServiceMover{}
	default:
		return &RandomMover{}
	}
}
