package game

type Position struct{
    Row int `json:"row"`
    Col int `json:"col"`
}

type Board [][] string

func NewBoard(size int)Board{
    board:=make(Board,size)
    for i:=range board{
        board[i]=make([]string,size)
    }
    return board
}
