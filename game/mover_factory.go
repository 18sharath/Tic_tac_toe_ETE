package game

func NewBotMover(difficulty Difficulty) Mover{
	switch difficulty{
	case DifficultyEasy:
		return &RandomeMover{}
	case DifficultyMedium:
		return  &DefensiveMover{}
	case DifficultyHard:
		return &OffensiveMover{}
	default:
		return &RandomeMover{}
	}
}