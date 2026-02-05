package game

// Snake holds the state of the snake on the board.
type Snake struct {
	Body      []Position
	Direction Direction
	Growing   bool
}

// NewSnake creates a snake starting at the given position, heading right.
func NewSnake(start Position) *Snake {
	return &Snake{
		Body: []Position{
			start,
			{X: start.X - 1, Y: start.Y},
			{X: start.X - 2, Y: start.Y},
		},
		Direction: Right,
	}
}

// Head returns the current head position.
func (s *Snake) Head() Position {
	return s.Body[0]
}

// Move advances the snake one step in its current direction.
// If Growing is true, the tail is not removed (the snake gets longer).
func (s *Snake) Move() {
	head := s.Head()
	var next Position

	switch s.Direction {
	case Up:
		next = Position{X: head.X, Y: head.Y - 1}
	case Down:
		next = Position{X: head.X, Y: head.Y + 1}
	case Left:
		next = Position{X: head.X - 1, Y: head.Y}
	case Right:
		next = Position{X: head.X + 1, Y: head.Y}
	}

	s.Body = append([]Position{next}, s.Body...)

	if s.Growing {
		s.Growing = false
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

// Grow marks the snake to grow on its next move.
func (s *Snake) Grow() {
	s.Growing = true
}

// SetDirection changes direction, preventing 180-degree reversal.
func (s *Snake) SetDirection(d Direction) {
	opposites := map[Direction]Direction{
		Up: Down, Down: Up, Left: Right, Right: Left,
	}
	if opposites[d] != s.Direction {
		s.Direction = d
	}
}

// CollidesWithSelf returns true if the head overlaps any body segment.
func (s *Snake) CollidesWithSelf() bool {
	head := s.Head()
	for _, seg := range s.Body[1:] {
		if seg == head {
			return true
		}
	}
	return false
}
