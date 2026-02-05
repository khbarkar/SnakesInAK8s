package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette and styles inspired by OpenClaw TUI.
// Warm dark background with golden accents.
type Theme struct {
	// Colors
	Background  lipgloss.Color
	Foreground  lipgloss.Color
	Accent      lipgloss.Color
	AccentSoft  lipgloss.Color
	Dim         lipgloss.Color
	Border      lipgloss.Color
	Error       lipgloss.Color
	Success     lipgloss.Color
	SnakeHead   lipgloss.Color
	SnakeBody   lipgloss.Color
	PodColor    lipgloss.Color

	// Derived styles
	HeaderStyle  lipgloss.Style
	FooterStyle  lipgloss.Style
	BoardStyle   lipgloss.Style
	ScoreStyle   lipgloss.Style
	StatusStyle  lipgloss.Style
	KillLogStyle lipgloss.Style
}

// DefaultTheme returns the OpenClaw-inspired color scheme.
func DefaultTheme() Theme {
	t := Theme{
		Background: lipgloss.Color("#2B2F36"),
		Foreground: lipgloss.Color("#E8E3D5"),
		Accent:     lipgloss.Color("#F6C453"),
		AccentSoft: lipgloss.Color("#F2A65A"),
		Dim:        lipgloss.Color("#7B7F87"),
		Border:     lipgloss.Color("#3C414B"),
		Error:      lipgloss.Color("#F97066"),
		Success:    lipgloss.Color("#7DD3A5"),
		SnakeHead:  lipgloss.Color("#F6C453"),
		SnakeBody:  lipgloss.Color("#F2A65A"),
		PodColor:   lipgloss.Color("#7DD3A5"),
	}

	t.HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent).
		Background(t.Background).
		Padding(0, 1)

	t.FooterStyle = lipgloss.NewStyle().
		Foreground(t.Dim).
		Background(t.Background).
		Padding(0, 1)

	t.BoardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(0)

	t.ScoreStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent)

	t.StatusStyle = lipgloss.NewStyle().
		Foreground(t.Dim).
		Italic(true)

	t.KillLogStyle = lipgloss.NewStyle().
		Foreground(t.Error)

	return t
}
