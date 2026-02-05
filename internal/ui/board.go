package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kristinb/snakeinak8/internal/game"
)

// Cell characters for rendering.
const (
	CellEmpty     = " "
	CellSnakeHead = "@"
	CellSnakeBody = "#"
	CellPod       = "*"
	CellWall      = "."
)

// RenderBoard draws the game board as a string.
func RenderBoard(theme Theme, g *game.Game) string {
	// Build a 2D grid
	grid := make([][]string, g.Board.Height)
	for y := range grid {
		grid[y] = make([]string, g.Board.Width)
		for x := range grid[y] {
			grid[y][x] = CellEmpty
		}
	}

	// Place pods
	podStyle := lipgloss.NewStyle().Foreground(theme.PodColor).Bold(true)
	for _, pod := range g.Pods {
		if inBounds(pod.Pos, g.Board) {
			grid[pod.Pos.Y][pod.Pos.X] = podStyle.Render(CellPod)
		}
	}

	// Place snake body
	bodyStyle := lipgloss.NewStyle().Foreground(theme.SnakeBody)
	for _, seg := range g.Snake.Body[1:] {
		if inBounds(seg, g.Board) {
			grid[seg.Y][seg.X] = bodyStyle.Render(CellSnakeBody)
		}
	}

	// Place snake head
	headStyle := lipgloss.NewStyle().Foreground(theme.SnakeHead).Bold(true)
	head := g.Snake.Head()
	if inBounds(head, g.Board) {
		grid[head.Y][head.X] = headStyle.Render(CellSnakeHead)
	}

	// Render rows
	var rows []string
	for _, row := range grid {
		rows = append(rows, strings.Join(row, ""))
	}
	board := strings.Join(rows, "\n")

	return theme.BoardStyle.Render(board)
}

func inBounds(p game.Position, b *game.Board) bool {
	return p.X >= 0 && p.X < b.Width && p.Y >= 0 && p.Y < b.Height
}
