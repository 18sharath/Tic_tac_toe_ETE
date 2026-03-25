package game

import "math/rand"

func (g *Game) BotMove() {
    switch g.Difficulty {
    case 1:
        g.randomMove()
    case 2:
        g.defensiveMove()
    case 3:
        g.offensiveMove()
    default:
        g.randomMove()
    }
}

func (g *Game) randomMove() {
    for {
        r := rand.Intn(3)
        c := rand.Intn(3)
        if g.Board[r][c] == "" {
            g.Board[r][c] = g.Player
            return
        }
    }
}

func (g *Game) defensiveMove() {
    opponent := g.OtherPlayer()

    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if g.Board[i][j] == "" {
                g.Board[i][j] = opponent
                if g.CheckWinner() == opponent {
                    g.Board[i][j] = g.Player
                    return
                }
                g.Board[i][j] = ""
            }
        }
    }
    g.randomMove()
}

func (g *Game) offensiveMove() {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if g.Board[i][j] == "" {
                g.Board[i][j] = g.Player
                if g.CheckWinner() == g.Player {
                    return
                }
                g.Board[i][j] = ""
            }
        }
    }
    g.defensiveMove()
}