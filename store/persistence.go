package store

import (
	"encoding/json"
	"os"
)

func LoadGames()error{
	file , err:=os.Open(dataFile)
	if err!=nil{
		if os.IsNotExist(err){
			return nil
		}
		return err
	}
	defer file.Close()

	Mutex.Lock()
	defer Mutex.Unlock()
	return json.NewDecoder(file).Decode(&Games)

}

func SaveGame() error{
	file,err:=os.Create(dataFile)
	if err!=nil{
		return  err
	}
	defer file.Close()
	Mutex.RLock()
	defer Mutex.RUnlock()
	encoder:=json.NewEncoder(file)
	encoder.SetIndent("","  ")
	return  encoder.Encode(Games)
}