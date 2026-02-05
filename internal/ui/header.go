package ui

import "github.com/charmbracelet/lipgloss"

// RenderHeader draws the top bar with game title and cluster info.
func RenderHeader(theme Theme, width int, clusterName string) string {
	title := theme.HeaderStyle.Render("snakeinak8")

	info := lipgloss.NewStyle().
		Foreground(theme.Dim).
		Render("cluster: " + clusterName)

	gap := width - lipgloss.Width(title) - lipgloss.Width(info)
	if gap < 1 {
		gap = 1
	}
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, title, spacer, info)

	return lipgloss.NewStyle().
		Background(theme.Background).
		Width(width).
		Render(bar)
}
