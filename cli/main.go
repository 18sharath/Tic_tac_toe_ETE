package main

import (
	"flag"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

var menuOptions = []string{
	"Human vs Human",
	"Human vs Bot",
	"Bot vs Bot",
	"Quit",
}

var difficultyOptions = []string{
	"Easy",
	"Medium",
	"Hard",
}

func initialModel() model {
	return model{
		screen: menuScreen,
	}
}

func (m model) Init() tea.Cmd { return nil }

func main() {
	var baseURLFlag = flag.String("base-url", "http://localhost:8080", "backend API base url")
	flag.Parse()

	baseURL = *baseURLFlag

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
