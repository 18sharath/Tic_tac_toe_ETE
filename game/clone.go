package game

// Clone creates a deep copy of the Game suitable for returning.
func (g *Game) Clone() *Game {
    if g == nil {
        return nil
    }

    var boardCopy Board
    if g.Board != nil {
        boardCopy = make(Board, len(g.Board))
        for i := range g.Board {
            boardCopy[i] = make([]string, len(g.Board[i]))
            copy(boardCopy[i], g.Board[i])
        }
    }

    return &Game{
        ID:         g.ID,
        Board:      boardCopy,
        PlayerX:    g.PlayerX,
        PlayerO:    g.PlayerO,
        Turn:       g.Turn,
        Winner:     g.Winner,
        Draw:       g.Draw,
        Mode:       g.Mode,
        Difficulty: g.Difficulty,
    }
}
