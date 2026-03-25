package game

func (g *Game) initBoard() {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            g.Board[i][j] = ""
        }
    }
}

func (g *Game) CheckWinner() string {
    b := g.Board

    for i := 0; i < 3; i++ {
        if b[i][0] != "" && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
            return b[i][0]
        }
        if b[0][i] != "" && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
            return b[0][i]
        }
    }

    if b[0][0] != "" && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
        return b[0][0]
    }

    if b[0][2] != "" && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
        return b[0][2]
    }

    return ""
}

func (g *Game) CheckDraw() bool {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if g.Board[i][j] == "" {
                return false
            }
        }
    }
    return true
}