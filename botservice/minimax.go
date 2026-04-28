package main

import "math"
const (
	scoreWin  = 10
	scoreLose = -10
	scoreInf  = math.MaxInt
)

func checkWinner(board [][]string) string {
	size := len(board)

	for i := 0; i < size; i++ {
		if board[i][0] != "" {
			win := true
			for j := 1; j < size; j++ {
				if board[i][j] != board[i][0] {
					win = false
					break
				}
			}
			if win {
				return board[i][0]
			}
		}
	}

	for j := 0; j < size; j++ {
		if board[0][j] != "" {
			win := true
			for i := 1; i < size; i++ {
				if board[i][j] != board[0][j] {
					win = false
					break
				}
			}
			if win {
				return board[0][j]
			}
		}
	}

	if board[0][0] != "" {
		win := true
		for i := 1; i < size; i++ {
			if board[i][i] != board[0][0] {
				win = false
				break
			}
		}
		if win {
			return board[0][0]
		}
	}

	if board[0][size-1] != "" {
		win := true
		for i := 1; i < size; i++ {
			if board[i][size-1-i] != board[0][size-1] {
				win = false
				break
			}
		}
		if win {
			return board[0][size-1]
		}
	}

	return ""
}

func isBoardFull(board [][]string) bool {
	for i := range board {
		for j := range board[i] {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

// minimax uses alpha-beta pruning to reduce unnecessary evaluations.
func minimax(board [][]string, isMaximizing bool, alpha, beta int) int {
	if score, done := terminalScore(board); done {
		return score
	}

	if isMaximizing {
		return maximizingScore(board, alpha, beta)
	}

	return minimizingScore(board, alpha, beta)
}

func terminalScore(board [][]string) (int, bool) {
	winner := checkWinner(board)
	if winner == "O" {
		return scoreWin, true
	}
	if winner == "X" {
		return scoreLose, true
	}
	if isBoardFull(board) {
		return 0, true
	}
	return 0, false
}

func maximizingScore(board [][]string, alpha, beta int) int {
	best := -scoreInf
	for i := range board {
		for j := range board[i] {
			if board[i][j] == "" {
				board[i][j] = "O"
				score := minimax(board, false, alpha, beta)
				board[i][j] = ""
				if score > best {
					best = score
				}
				if best > alpha {
					alpha = best
				}
				if beta <= alpha {
					return best
				}
			}
		}
	}
	return best
}

func minimizingScore(board [][]string, alpha, beta int) int {
	best := scoreInf
	for i := range board {
		for j := range board[i] {
			if board[i][j] == "" {
				board[i][j] = "X"
				score := minimax(board, true, alpha, beta)
				board[i][j] = ""
				if score < best {
					best = score
				}
				if best < beta {
					beta = best
				}
				if beta <= alpha {
					return best
				}
			}
		}
	}
	return best
}

func bestMove(board [][]string) moveResponse {
	best := -scoreInf
	move := moveResponse{Row: -1, Col: -1}

	for i := range board {
		for j := range board[i] {
			if board[i][j] == "" {
				board[i][j] = "O"
				score := minimax(board, false, -scoreInf, scoreInf)
				board[i][j] = ""
				if score > best {
					best = score
					move = moveResponse{Row: i, Col: j}
				}
			}
		}
	}
	return move
}
