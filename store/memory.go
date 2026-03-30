package store

import (
	"tic_tac_toe/game"
	
	"sync"
	
	)

var (
	Games = make(map[string]*game.Game)
	Mutex sync.RWMutex
	dataFile="data/games.json"
)

