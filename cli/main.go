package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
