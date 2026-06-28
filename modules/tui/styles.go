package tui

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#3A3A3A"}
	highlight = lipgloss.AdaptiveColor{Light: "#4ecdc4", Dark: "#4ecdc4"}
	special   = lipgloss.AdaptiveColor{Light: "#ff6b6b", Dark: "#ff6b6b"}
	green     = lipgloss.AdaptiveColor{Light: "#2ecc71", Dark: "#2ecc71"}
	yellow    = lipgloss.AdaptiveColor{Light: "#f1c40f", Dark: "#f1c40f"}
	red       = lipgloss.AdaptiveColor{Light: "#e74c3c", Dark: "#e74c3c"}
	white     = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
	black     = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#000000"}

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Padding(0, 1).
			Width(60)

	statusStyle = lipgloss.NewStyle().
			Foreground(white).
			Padding(0, 1)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(subtle)

	activeTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(highlight).
			Bold(true).
			Underline(true)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Width(20)

	valueStyle = lipgloss.NewStyle().
			Foreground(white)

	labelStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Width(12)

	connectedStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	disconnectedStyle = lipgloss.NewStyle().
				Foreground(red).
				Bold(true)

	peerNameStyle = lipgloss.NewStyle().
			Foreground(white).
			Width(15)

	peerIPStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Width(15)

	peerTrafficStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Width(20)

	logTimeStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Width(8)

	logTextStyle = lipgloss.NewStyle().
			Foreground(white)

	buttonStyle = lipgloss.NewStyle().
			Foreground(black).
			Background(highlight).
			Padding(0, 2).
			Bold(true)

	buttonDisabledStyle = lipgloss.NewStyle().
				Foreground(subtle).
				Background(lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#2A2A2A"}).
				Padding(0, 2)

	commandPromptStyle = lipgloss.NewStyle().
				Foreground(green).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(0, 1)

	commandInputStyle = lipgloss.NewStyle().
				Foreground(white).
				Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#1A1A1A"})

	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(1)

	dashStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(subtle).
			Padding(0, 1).
			Width(30)
)
