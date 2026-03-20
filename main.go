package main

import (
	"fmt"
	"math/rand"
	"time"
)	

var board [3][3]string
var mode int

func initBoard(){
	for i:=0;i<3;i++{
		for j:=0;j<3;j++{
			board[i][j]=""
		}
	}
}

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

func Botmove(Bot string){
	for{
		row:=rand.Intn(3)
		col:=rand.Intn(3)

		if board[row][col]==""{
			fmt.Printf("Bot %s played:[%v %v]\n",Bot,row,col)
			board[row][col]=Bot
			break
		}

	}
}


func PlayerMove(player string){
	var row,col int
	for{
		fmt.Printf("Player %s - Enter row and column:",player)
		fmt.Scan(&row, &col)

		if row<0 || row >=3 || col<0 || col>=3 || board[row][col]!=""{
			fmt.Println("Invalid move, try again!")
			continue
		}
		board[row][col] = player
			break
	}
}


func main(){	
	player:="X"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	initBoard()

	fmt.Println("Choose Mode:")
	fmt.Println("1. Player vs Player")
	fmt.Println("2. Player vs Bot")
	fmt.Println("3. Bot vs Bot")
	fmt.Scan(&mode)

	for{

		printBoard()
	
		if mode==1|| (mode==2 && player=="X"){
			PlayerMove(player)
		}else{
			Botmove(player)
		}

		winner:=checkwinner()
		if winner != "" {
			printBoard()
			if mode==2&&winner=="O"{
				fmt.Println("Bot wins")
			}else{
				fmt.Printf("\n Player %s wins!\n", winner)
			}
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

