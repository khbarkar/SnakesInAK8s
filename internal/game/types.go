package game

// Direction represents the snake's movement direction.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Position represents a coordinate on the game board.
type Position struct {
	X int
	Y int
}

// Pod represents a Kubernetes pod displayed on the board.
type Pod struct {
	Pos       Position
	Name      string
	Namespace string
}
