package store

import "tic_tac_toe/game"

type Store interface{
	Create(game *game.Game)error
	Get(id string)(*game.Game,bool)
	Delete(id string)error
}