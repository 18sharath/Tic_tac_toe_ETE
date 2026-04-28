package game

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var botserviceURL string

// SetBotServiceURL configures the external bot service endpoint.
func SetBotServiceURL(url string) {
	botserviceURL=url
}

type moveRequest struct{
	Board Board `json:"board"`
	Player string `json:"player"`
}

type moveResponse struct{
	Row int `json:"row"`
	Col int `json:"col"`
}

// ServiceMover gets bot moves from external bot services.
type ServiceMover struct{}

// Move calls external bot moves and returns next move.
func (s *ServiceMover) Move(board Board,  player string) (Position,error) {
	reqBody := moveRequest{
		Board: board,
		Player: player,
	}

	jsonData, err := json.Marshal(reqBody)

	if err!=nil{
		return Position{}, err
	}

	resp, err := http.Post(
		botserviceURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err!=nil{
		return Position{}, err
	}

	defer func() {
			if cerr := resp.Body.Close(); cerr != nil && err != nil {
				err = cerr
			}
		}()

	var result moveResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err!=nil{
		return  Position{}, err
	}

	return Position(result),nil
}