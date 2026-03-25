package game
import "fmt"
type Game struct {
    ID         string      `json:"id"`
    Board      [3][3]string`json:"board"`
    Player     string      `json:"player"`
    Mode       int         `json:"mode"`
    Difficulty int         `json:"difficulty"`
    Winner     string      `json:"winner"`
    Draw       bool        `json:"draw"`
}


func NewGame(id string, mode, difficulty int) *Game {
    g := &Game{
        ID:         id,
        Player:     "X",
        Mode:       mode,
        Difficulty: difficulty,
    }
    g.initBoard()
    return g
}

func (g *Game) TogglePlayer() {
    if g.Player == "X" {
        g.Player = "O"
    } else {
        g.Player = "X"
    }
}

func (g *Game) OtherPlayer() string {
    if g.Player == "X" {
        return "O"
    }
    return "X"
}

func (g *Game) MakeMove(row, col int) error{
    
 if row < 0 || row > 2 || col < 0 || col > 2 {
        return fmt.Errorf("invalid board position")
    }

    if g.Board[row][col] != "" {
        return fmt.Errorf("cell already occupied")
    }
	
    g.Board[row][col] = g.Player
    return nil

}


func (g *Game) IsGameOver() bool {
    return g.Winner != "" || g.Draw
}
