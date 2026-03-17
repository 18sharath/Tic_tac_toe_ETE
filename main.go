package main

import "fmt"

var board [3][3]string


func printBoard() {
	fmt.Println("Current Board")

	for i:=0;i<3;i++{
		for j:=0;j<3;j++{
			if board[i][j]==""{
				fmt.Print("-")
			}else{
				fmt.Print(board[i][j]+" ")
			}
		}
		fmt.Println()
	}
}

func checkwinner() string{
	
	for i:=0;i<3;i++{
		// row
		if board[i][0]!="" && board[i][0]==board[i][1] &&board[i][1]==board[i][2]{
			return board[i][0]
		}
		// col
		if board[0][i]!="" && board[0][i]==board[1][i] &&board[1][i]==board[2][i]{
			return board[0][i]
		}
	}

	if board[0][0] != "" && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0]
	}
	if board[0][2] != "" && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2]
	}

	return ""
}

func drawConditionCheck() bool{
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func main(){	
	player:="X"

	for{

		printBoard()
		var row,col int

		fmt.Printf("Player %s Enter row and column:",player)
		fmt.Scan(&row, &col)

		if row<0 || col <0 || row >2 || col>2 || board[row][col]!=""{
			fmt.Println("Invalid move! Try again ")
			continue
		}

		board[row][col]=player

		winner:=checkwinner()
		if winner != "" {
			printBoard()
			fmt.Printf("\n Player %s wins!\n", winner)
			break
		}
		if drawConditionCheck() {
			printBoard()
			fmt.Println("\nIt's a draw!")
			break
		}

		if player == "X" {
			player = "O"
		} else {
			player = "X"
		}


	}
}


