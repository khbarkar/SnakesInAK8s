package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kristinb/snakeinak8/internal/game"
	"github.com/kristinb/snakeinak8/internal/k8s"
)

const (
	defaultTickRate = 150 * time.Millisecond
	boardWidth      = 40
	boardHeight     = 20
)

// tickMsg fires on every game tick.
type tickMsg time.Time

// podPlacedMsg signals a new pod was fetched from k8s and should appear.
type podPlacedMsg struct {
	Name      string
	Namespace string
	Err       error
}

// podKilledMsg signals a pod was deleted from the cluster.
type podKilledMsg struct {
	Name      string
	Namespace string
	Err       error
}

// GameModel is the top-level Bubble Tea model for the game.
type GameModel struct {
	game        *game.Game
	theme       Theme
	killLog     []string
	clusterName string
	namespace   string
	k8sClient   *k8s.Client
	width       int
	height      int
	tickRate    time.Duration
	fetching    bool // true while a pod fetch is in flight
}

// NewGameModel creates the game model with a connected k8s client.
func NewGameModel(client *k8s.Client, namespace string, theme Theme) GameModel {
	clusterName := "unknown"
	if client != nil {
		clusterName = client.ClusterName()
	}

	return GameModel{
		game:        game.New(boardWidth, boardHeight),
		theme:       theme,
		killLog:     []string{},
		clusterName: clusterName,
		namespace:   namespace,
		k8sClient:   client,
		tickRate:    defaultTickRate,
	}
}

// Init starts the tick loop and fetches initial pods.
func (m GameModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(m.tickRate),
		fetchPodCmd(m.k8sClient),
		fetchPodCmd(m.k8sClient),
		fetchPodCmd(m.k8sClient),
	)
}

// Update handles messages (keys, ticks, pod events).
func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "w":
			m.game.Snake.SetDirection(game.Up)
		case "down", "s":
			m.game.Snake.SetDirection(game.Down)
		case "left", "a":
			m.game.Snake.SetDirection(game.Left)
		case "right", "d":
			m.game.Snake.SetDirection(game.Right)
		case " ":
			m.game.TogglePause()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		eaten := m.game.Tick()
		var cmds []tea.Cmd

		for _, pod := range eaten {
			m.killLog = append(m.killLog, pod.Namespace+"/"+pod.Name)
			cmds = append(cmds, killPodCmd(m.k8sClient, pod.Name, pod.Namespace))
		}

		// Replenish pods on the board
		if len(m.game.Pods) < m.game.MaxPods && !m.fetching {
			m.fetching = true
			cmds = append(cmds, fetchPodCmd(m.k8sClient))
		}

		cmds = append(cmds, tickCmd(m.tickRate))
		return m, tea.Batch(cmds...)

	case podPlacedMsg:
		m.fetching = false
		if msg.Err == nil && msg.Name != "" {
			m.game.PlacePod(msg.Name, msg.Namespace)
		}

	case podKilledMsg:
		if msg.Err != nil {
			m.killLog = append(m.killLog, "FAILED: "+msg.Namespace+"/"+msg.Name+" -- "+msg.Err.Error())
		}
	}

	return m, nil
}

// View renders the full UI.
func (m GameModel) View() string {
	if m.width == 0 {
		return "initializing..."
	}

	stateLabel := "running"
	switch m.game.State {
	case game.StatePaused:
		stateLabel = "paused"
	case game.StateOver:
		stateLabel = "GAME OVER"
	}

	header := RenderHeader(m.theme, m.width, m.clusterName)
	board := RenderBoard(m.theme, m.game)
	footer := RenderFooter(m.theme, m.width, m.game.Score, m.game.KillCount, stateLabel)

	// Kill log: show last 5 kills
	var killLines string
	start := len(m.killLog) - 5
	if start < 0 {
		start = 0
	}
	for _, entry := range m.killLog[start:] {
		killLines += m.theme.KillLogStyle.Render("  killed: "+entry) + "\n"
	}

	nsLabel := "all"
	if m.namespace != "" {
		nsLabel = m.namespace
	}
	controls := m.theme.FooterStyle.Render("[wasd/arrows] move  [space] pause  [q] quit  ns:" + nsLabel)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		board,
		"",
		killLines,
		footer,
		controls,
	)
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchPodCmd(client *k8s.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return podPlacedMsg{Err: nil}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pod, err := client.RandomPod(ctx)
		if err != nil {
			return podPlacedMsg{Err: err}
		}
		if pod == nil {
			return podPlacedMsg{}
		}
		return podPlacedMsg{Name: pod.Name, Namespace: pod.Namespace}
	}
}

func killPodCmd(client *k8s.Client, name, namespace string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return podKilledMsg{Name: name, Namespace: namespace}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := client.KillPod(ctx, name, namespace)
		return podKilledMsg{Name: name, Namespace: namespace, Err: err}
	}
}
