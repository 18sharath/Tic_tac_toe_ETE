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

<<<<<<< HEAD
<<<<<<< HEAD
// NewFileStore helps to assign file to dataFile
=======
//NewFIleStore helps to assign file to dataFile
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
// NewFileStore helps to assign file to dataFile
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
func NewFileStore(file string) *FileStore {
	return &FileStore{
		dataFile: file,
	}
}

// Create function  helps to create new file and save the game to file
func (f *FileStore) Create(g *game.Game) (err error) {
<<<<<<< HEAD
<<<<<<< HEAD
	fileName := g.ID + ".json"
	filePath := filepath.Join(f.dataFile, fileName)
=======

	filePath := f.dataFile + "/" + g.ID + ".json"

>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
	fileName := g.ID + ".json"
	filePath := filepath.Join(f.dataFile, fileName)
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
<<<<<<< HEAD
<<<<<<< HEAD
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = json.NewEncoder(file).Encode(g)
=======
	defer func(){
		if cerr:=file.Close();cerr!=nil && err==nil{
			err=cerr
		}
	}()

	err=json.NewEncoder(file).Encode(g)
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = json.NewEncoder(file).Encode(g)
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
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

<<<<<<< HEAD
<<<<<<< HEAD
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close the file %s: %v", filePath, err)
=======
	defer func(){
		if err:=file.Close();err!=nil{
			log.Printf("Failed to close the file %s: %v",filePath,err)
>>>>>>> ed86cc4 (This change will update the codebase based on the comments on PR)
=======
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close the file %s: %v", filePath, err)
>>>>>>> eb32414 (This change will make more production-ready based on golangci-lint)
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
