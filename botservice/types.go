package main

type moveRequest struct {
	Board  [][]string `json:"board"`
	Player string     `json:"player"`
}

type moveResponse struct {
	Row int `json:"row"`
	Col int `json:"col"`
}
