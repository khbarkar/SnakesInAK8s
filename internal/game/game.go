package game

// State represents the current phase of the game.
type State int

const (
	StateRunning State = iota
	StatePaused
	StateOver
)

// Game ties together the snake, board, and pod targets.
type Game struct {
	Snake      *Snake
	Board      *Board
	Pods       []Pod
	State      State
	Score      int
	KillCount  int
	MaxPods    int // maximum pods visible on board at once
}

// New creates a new game with default settings.
// The board dimensions are passed in so the UI can control sizing.
func New(boardWidth, boardHeight int) *Game {
	board := NewBoard(boardWidth, boardHeight)
	start := Position{X: boardWidth / 4, Y: boardHeight / 2}
	snake := NewSnake(start)

	return &Game{
		Snake:   snake,
		Board:   board,
		Pods:    []Pod{},
		State:   StateRunning,
		MaxPods: 3,
	}
}

// Tick advances the game by one frame. Returns a list of pod names that
// were eaten (and should be killed in k8s).
func (g *Game) Tick() []Pod {
	if g.State != StateRunning {
		return nil
	}

	g.Snake.Move()
	head := g.Snake.Head()

	// TODO: wall collision -- game over
	if g.Board.IsOutOfBounds(head) {
		g.State = StateOver
		return nil
	}

	// TODO: self collision -- game over
	if g.Snake.CollidesWithSelf() {
		g.State = StateOver
		return nil
	}

	// Check if we ate a pod
	var eaten []Pod
	remaining := make([]Pod, 0, len(g.Pods))
	for _, pod := range g.Pods {
		if pod.Pos == head {
			eaten = append(eaten, pod)
			g.Snake.Grow()
			g.Score++
			g.KillCount++
		} else {
			remaining = append(remaining, pod)
		}
	}
	g.Pods = remaining

	return eaten
}

// PlacePod adds a pod to the board at a random free position.
// Returns true if the pod was placed, false if the board is full.
func (g *Game) PlacePod(name, namespace string) bool {
	if len(g.Pods) >= g.MaxPods {
		return false
	}

	occupied := g.Snake.Body
	for _, p := range g.Pods {
		occupied = append(occupied, p.Pos)
	}

	pos := g.Board.RandomPosition(occupied)
	g.Pods = append(g.Pods, Pod{
		Pos:       pos,
		Name:      name,
		Namespace: namespace,
	})
	return true
}

// TogglePause pauses or resumes the game.
func (g *Game) TogglePause() {
	switch g.State {
	case StateRunning:
		g.State = StatePaused
	case StatePaused:
		g.State = StateRunning
	}
}
