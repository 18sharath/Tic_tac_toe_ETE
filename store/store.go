// Package store contains interfaces that define game storage behavior.
package store

import "tic_tac_toe/game"

// GameStore defines the behavior required for storing and managing games.
type GameStore interface {
	Create(game *game.Game) error
	Get(id string) (*game.Game, bool)
	Delete(id string) error
}
