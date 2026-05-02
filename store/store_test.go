package store

import (
	"os"
	"testing"
	"tic_tac_toe/game"

	"github.com/stretchr/testify/assert"
)

func createTestGame(id string) *game.Game {
	return game.NewGame(
		id,
		3,
		game.ModeHumanVsHuman,
		game.DifficultyEasy,
		nil,
		nil,
	)
}

func TestMemoryStoreCreateAndGet(t *testing.T) {
	store := NewMemoryStore()

	g := createTestGame("game-1")

	err := store.Create(g)

	assert.NoError(t, err)

	result, ok := store.Get("game-1")

	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.Equal(t, "game-1", result.ID)
}

func TestMemoryStoreReturnsCopies(t *testing.T) {
	store := NewMemoryStore()

	g := createTestGame("copy-game")
	g.Board[0][0] = "X"

	err := store.Create(g)

	assert.NoError(t, err)

	result, ok := store.Get("copy-game")

	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.Equal(t, "X", result.Board[0][0])

	result.Board[0][0] = "O"

	result2, ok := store.Get("copy-game")

	assert.True(t, ok)
	assert.Equal(t, "X", result2.Board[0][0])
	assert.NotSame(t, result, result2)
}

func TestMemoryStoreGetNotFound(t *testing.T) {
	store := NewMemoryStore()

	result, ok := store.Get("missing-id")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestMemoryStoreDelete(t *testing.T) {
	store := NewMemoryStore()

	g := createTestGame("delete-id")
	_ = store.Create(g)

	err := store.Delete("delete-id")

	assert.NoError(t, err)

	result, ok := store.Get("delete-id")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestFileStoreCreateAndGet(t *testing.T) {
	tempDir := t.TempDir()

	store := NewFileStore(tempDir)

	g := createTestGame("file-game-1")

	err := store.Create(g)

	assert.NoError(t, err)

	result, ok := store.Get("file-game-1")

	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.Equal(t, "file-game-1", result.ID)
	assert.Equal(t, "X", result.Turn)
}

func TestFileStoreGetNotFound(t *testing.T) {
	tempDir := t.TempDir()

	store := NewFileStore(tempDir)

	result, ok := store.Get("not-found")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestFileStoreDelete(t *testing.T) {
	tempDir := t.TempDir()

	store := NewFileStore(tempDir)

	g := createTestGame("delete-file")
	_ = store.Create(g)

	err := store.Delete("delete-file")

	assert.NoError(t, err)

	result, ok := store.Get("delete-file")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestFileStoreDeleteMissingFile(t *testing.T) {
	tempDir := t.TempDir()

	store := NewFileStore(tempDir)

	err := store.Delete("missing-file")

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}