package store

import (
	"tic_tac_toe/game"
	"sync"
)

type MemoryStore struct {
	games map[string]*game.Game
	mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games: make(map[string]*game.Game),
	}
}

func (m *MemoryStore) Create(g *game.Game) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.games[g.ID] = g
	return nil
}

func (m *MemoryStore) Get(id string) (*game.Game, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	g, ok := m.games[id]
	return g, ok
}

func (m *MemoryStore) Delete(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.games, id)
	return nil
}


var (
	Games    = make(map[string]*game.Game)
	Mutex    sync.RWMutex
	dataFile = "data/games.json"
)
