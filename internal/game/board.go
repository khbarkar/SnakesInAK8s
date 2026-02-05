package game

import (
	"math/rand"
)

// Board holds the dimensions and state of the game area.
type Board struct {
	Width  int
	Height int
}

// NewBoard creates a board with the given dimensions.
func NewBoard(width, height int) *Board {
	return &Board{Width: width, Height: height}
}

// IsOutOfBounds returns true if the position is outside the board.
func (b *Board) IsOutOfBounds(p Position) bool {
	return p.X < 0 || p.X >= b.Width || p.Y < 0 || p.Y >= b.Height
}

// RandomPosition returns a random position within the board that does not
// overlap with any of the excluded positions.
func (b *Board) RandomPosition(excluded []Position) Position {
	excludeSet := make(map[Position]bool, len(excluded))
	for _, p := range excluded {
		excludeSet[p] = true
	}

	for {
		p := Position{
			X: rand.Intn(b.Width),
			Y: rand.Intn(b.Height),
		}
		if !excludeSet[p] {
			return p
		}
	}
}
