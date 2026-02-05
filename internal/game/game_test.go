package game

import "testing"

func TestNewSnake(t *testing.T) {
	s := NewSnake(Position{X: 5, Y: 5})
	if len(s.Body) != 3 {
		t.Fatalf("expected snake length 3, got %d", len(s.Body))
	}
	if s.Head() != (Position{X: 5, Y: 5}) {
		t.Fatalf("expected head at (5,5), got %v", s.Head())
	}
}

func TestSnakeMove(t *testing.T) {
	s := NewSnake(Position{X: 5, Y: 5})
	s.Move()
	if s.Head() != (Position{X: 6, Y: 5}) {
		t.Fatalf("expected head at (6,5) after move right, got %v", s.Head())
	}
	if len(s.Body) != 3 {
		t.Fatalf("expected snake length 3, got %d", len(s.Body))
	}
}

func TestSnakeGrow(t *testing.T) {
	s := NewSnake(Position{X: 5, Y: 5})
	s.Grow()
	s.Move()
	if len(s.Body) != 4 {
		t.Fatalf("expected snake length 4 after grow+move, got %d", len(s.Body))
	}
}

func TestSnakeReversal(t *testing.T) {
	s := NewSnake(Position{X: 5, Y: 5})
	s.SetDirection(Left)
	if s.Direction != Right {
		t.Fatal("snake should not reverse 180 degrees")
	}
}

func TestBoardOutOfBounds(t *testing.T) {
	b := NewBoard(10, 10)
	if b.IsOutOfBounds(Position{X: 0, Y: 0}) {
		t.Fatal("(0,0) should be in bounds")
	}
	if !b.IsOutOfBounds(Position{X: -1, Y: 0}) {
		t.Fatal("(-1,0) should be out of bounds")
	}
	if !b.IsOutOfBounds(Position{X: 10, Y: 5}) {
		t.Fatal("(10,5) should be out of bounds on 10-wide board")
	}
}

func TestGameTick(t *testing.T) {
	g := New(20, 20)
	g.PlacePod("test-pod", "default")

	if len(g.Pods) != 1 {
		t.Fatalf("expected 1 pod, got %d", len(g.Pods))
	}

	// Run ticks until game over or we eat a pod (limited iterations)
	for i := 0; i < 100; i++ {
		eaten := g.Tick()
		if len(eaten) > 0 {
			if eaten[0].Name != "test-pod" {
				t.Fatalf("expected eaten pod name 'test-pod', got %q", eaten[0].Name)
			}
			return
		}
		if g.State == StateOver {
			// Snake hit wall or self before reaching pod -- acceptable in test
			return
		}
	}
}

func TestGameOver(t *testing.T) {
	// Create a tiny board so the snake hits a wall quickly
	g := New(4, 4)
	g.Snake = NewSnake(Position{X: 2, Y: 2})

	for i := 0; i < 10; i++ {
		g.Tick()
	}
	if g.State != StateOver {
		t.Fatal("expected game over after snake runs off small board")
	}
}
