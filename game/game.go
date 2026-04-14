package game

import (
	"errors"
)

type Game struct {
	ID      string `json:"id"`
	Board   Board  `json:"board"`
	PlayerX Mover  `json:"-"`
	PlayerO Mover  `json:"-"`
	Turn    string `json:"turn"`
	Winner  string `json:"winner"`
	Draw    bool   `json:"draw"`
    Mode    Mode    `json:"mode"`
    Difficulty Difficulty `json:"difficulty"`
   
}

func NewGame(id string, size int, mode Mode, difficulty Difficulty, xMover Mover, oMover Mover) *Game {
	return &Game{
		ID:      id,
		Board:   NewBoard(size),
		PlayerX: xMover,
		PlayerO: oMover,
		Turn:    "X",
        Mode:       mode,
		Difficulty: difficulty, 
	}
}

func (g *Game) Maketurn() error {
	var mover Mover
	if g.Turn == "X" {
		mover = g.PlayerX
	} else {
		mover = g.PlayerO
	}

	if mover == nil {
		return errors.New("no mover available for current player")
	}

	pos, err := mover.Move(g.Board,g.Turn)
	if err != nil {
		return err
	}
	if g.Board[pos.Row][pos.Col] != "" {
		return errors.New("Invalid Move")
	}
	g.Board[pos.Row][pos.Col] = g.Turn
	g.toggleTurn()
	return nil

}

func (g *Game) toggleTurn() {
	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}
}

func (g *Game) Evaluate() {
	size := len(g.Board)

	//  Check rows
	for i := 0; i < size; i++ {
		first := g.Board[i][0] 
		if first == "" {
			continue
		}
		win := true
		for j := 1; j < size; j++ { 
			if g.Board[i][j] != first {
				win = false
				break
			}
		}
		if win {
			g.Winner = first
			return
		}
	}

	//  Check columns
	for j := 0; j < size; j++ {
		first := g.Board[0][j]
		if first == "" {
			continue
		}
		win := true
		for i := 1; i < size; i++ {
			if g.Board[i][j] != first {
				win = false
				break
			}
		}
		if win {
			g.Winner = first
			return
		}
	}

	//  Check main diagonal
	first := g.Board[0][0]
	if first != "" {
		win := true
		for i := 1; i < size; i++ {
			if g.Board[i][i] != first {
				win = false
				break
			}
		}
		if win {
			g.Winner = first
			return
		}
	}

	//  Check anti-diagonal
	first = g.Board[0][size-1]
	if first != "" {
		win := true
		for i := 1; i < size; i++ {
			if g.Board[i][size-1-i] != first {
				win = false
				break
			}
		}
		if win {
			g.Winner = first
			return
		}
	}

	//  Check draw
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if g.Board[i][j] == "" {
				return
			}
		}
	}

	// No empty cells and no winner
	g.Draw = true
}

func (g *Game) MakeMove(player string, row,col int) error{
    if player!= g.Turn{
        return errors.New("Not your turn")
    }
    size := len(g.Board)

    if row < 0 || row>= size || col < 0 || col>=size{
        return errors.New("invalid position")
    }

    if g.Board[row][col]!=""{
        return errors.New("Cell already occupied")
    }

    g.Board[row][col]=player

    g.toggleTurn()
    return nil
}

