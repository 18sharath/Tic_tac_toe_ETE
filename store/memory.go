package store

import (
	"sync"
	"tic_tac_toe/game"
)

var (
	// Games store all the games indexed by their unique ID.
	Games = make(map[string]*game.Game)

	// Mutex protects concurrent access to the in-memory game store.
	Mutex sync.RWMutex
)

// MemoryStore create map to store games in memory.
type MemoryStore struct {
	games map[string]*game.Game
	mutex sync.RWMutex
}

// NewMemoryStore creates and returns a new in-memory GameStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games: make(map[string]*game.Game),
	}
}

// Create helps to store new game in map
func (m *MemoryStore) Create(g *game.Game) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.games[g.ID] = g
	return nil
}

// Get helps to fetch the games from the map
func (m *MemoryStore) Get(id string) (*game.Game, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	g, ok := m.games[id]
	return g, ok
}

// Delete helps to remove the games from the map
func (m *MemoryStore) Delete(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.games, id)
	return nil
}
