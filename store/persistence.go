package store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"tic_tac_toe/game"
)

// FileStore contain the dataFIle path
type FileStore struct {
	dataFile string
}

// NewFileStore helps to assign file to dataFile
func NewFileStore(file string) *FileStore {
	return &FileStore{
		dataFile: file,
	}
}

// Create function  helps to create new file and save the game to file
func (f *FileStore) Create(g *game.Game) (err error) {
	fileName := g.ID + ".json"
	filePath := filepath.Join(f.dataFile, fileName)
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = json.NewEncoder(file).Encode(g)
	return err
}

// Get helps to open the stored game from the file
func (f *FileStore) Get(id string) (*game.Game, bool) {
	fileName := id + ".json"
	filePath := filepath.Join(f.dataFile, fileName)

	file, err := os.Open(filePath)

	if err != nil {
		log.Println("Looking for file:", filePath)
		return nil, false
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close the file %s: %v", filePath, err)
		}
	}()

	var g game.Game

	if err := json.NewDecoder(file).Decode(&g); err != nil {
		return nil, false
	}
	switch g.Mode {
	case game.ModeHumanVsHuman:
		g.PlayerX = nil
		g.PlayerO = nil

	case game.ModeHumanVsBot:
		g.PlayerX = nil
		g.PlayerO = game.NewBotMover(g.Difficulty)

	case game.ModeBotVsBot:
		g.PlayerX = game.NewBotMover(g.Difficulty)
		g.PlayerO = game.NewBotMover(g.Difficulty)
	}
	return &g, true
}

// Delete helps to delete stored game from the file
func (f *FileStore) Delete(id string) error {
	fileName := id + ".json"
	filePath := filepath.Join(f.dataFile, fileName)
	return os.Remove(filePath)
}
