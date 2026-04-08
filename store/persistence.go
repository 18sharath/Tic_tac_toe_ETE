package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tic_tac_toe/game"
)

type FileStore struct {
	dataFile string
}

func NewFileStore(file string) *FileStore {
	return &FileStore{
		dataFile: file,
	}
}

func (f *FileStore) Create(g *game.Game) error {

	filePath := f.dataFile + "/" + g.ID + ".json"

	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(g)
}

func (f *FileStore) Get(id string) (*game.Game, bool) {
	// filePath := f.dataFile + "/" + id + ".json"
	filePath:= filepath.Join(f.dataFile,id+".json")

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Looking for file:", filePath)
		return nil, false
	}
	defer file.Close()

	var g game.Game
	
	if err := json.NewDecoder(file).Decode(&g); err != nil {
		return nil, false
	}
	switch g.Mode{
	case game.ModeHumanVsHuman:
		g.PlayerX=nil
		g.PlayerO=nil
	
	case game.ModeHumanVsBot:
		g.PlayerX=	nil
		g.PlayerO= game.NewBotMover(g.Difficulty)

	case game.ModeBotVsBot:
		g.PlayerX=game.NewBotMover(g.Difficulty)
		g.PlayerO=game.NewBotMover(g.Difficulty)
	}
	return &g, true
}

func (f *FileStore) Delete(id string) error {
	filePath := f.dataFile + "/" + id + ".json"

	return os.Remove(filePath)
}
