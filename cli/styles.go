package main
import "github.com/charmbracelet/lipgloss"


var (
	xStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true) // green
	oStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true) // red
	cursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Bold(true)

	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

var (
	menuTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true).
			Padding(1, 2)

	menuItemStyle = lipgloss.NewStyle().
			Padding(0, 4)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("2")). // green background
			Bold(true).
			Padding(0, 4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			PaddingTop(1)
)
var popupStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 4).
	BorderForeground(lipgloss.Color("5")).
	Background(lipgloss.Color("0")).
	Align(lipgloss.Center)
