package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderFooter draws the bottom status bar with score and kill count.
func RenderFooter(theme Theme, width, score, kills int, state string) string {
	left := theme.ScoreStyle.Render(fmt.Sprintf("score: %d", score))
	mid := theme.KillLogStyle.Render(fmt.Sprintf("pods killed: %d", kills))
	right := theme.StatusStyle.Render(state)

	totalContent := lipgloss.Width(left) + lipgloss.Width(mid) + lipgloss.Width(right)
	gaps := width - totalContent
	if gaps < 2 {
		gaps = 2
	}
	spacer := lipgloss.NewStyle().Width(gaps / 2).Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, mid, spacer, right)

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(width).
		Render(bar)
}
