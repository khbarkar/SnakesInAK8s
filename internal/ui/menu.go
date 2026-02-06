package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kristinb/snakeinak8/internal/k8s"
)

// menuState tracks which screen the menu is on.
type menuState int

const (
	menuMain menuState = iota
	menuNamespace
	menuConnecting
	menuError
)

// namespacesLoadedMsg carries the list of namespaces from the cluster.
type namespacesLoadedMsg struct {
	namespaces []string
	err        error
}

// k8sConnectedMsg signals the k8s client was successfully created.
type k8sConnectedMsg struct {
	client *k8s.Client
	err    error
}

// MenuModel is the pre-game menu for configuring kubeconfig and namespace.
type MenuModel struct {
	theme          Theme
	kubeconfigPath string
	namespace      string // empty = all
	namespaces     []string
	cursor         int
	state          menuState
	errMsg         string
	k8sClient      *k8s.Client
	width          int
	height         int
	clusterName    string
}

// NewMenuModel creates the menu with the resolved kubeconfig path.
func NewMenuModel(kubeconfigPath string) MenuModel {
	return MenuModel{
		theme:          DefaultTheme(),
		kubeconfigPath: kubeconfigPath,
		namespace:      "",
		state:          menuConnecting,
	}
}

// NewMenuModelFromGame rebuilds the menu from a running game, preserving
// the k8s client, namespace selection, and terminal dimensions.
func NewMenuModelFromGame(g GameModel) MenuModel {
	return MenuModel{
		theme:          g.theme,
		kubeconfigPath: g.kubeconfig,
		namespace:      g.namespace,
		k8sClient:      g.k8sClient,
		clusterName:    g.clusterName,
		width:          g.width,
		height:         g.height,
		state:          menuMain,
		cursor:         0,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return connectK8sCmd(m.kubeconfigPath)
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state != menuNamespace {
				return m, tea.Quit
			}
		}

		switch m.state {
		case menuMain:
			return m.updateMain(msg)
		case menuNamespace:
			return m.updateNamespace(msg)
		case menuError:
			return m.updateError(msg)
		}

	case k8sConnectedMsg:
		if msg.err != nil {
			m.state = menuError
			m.errMsg = msg.err.Error()
			return m, nil
		}
		m.k8sClient = msg.client
		m.clusterName = msg.client.ClusterName()
		m.state = menuMain
		m.cursor = 0
		return m, nil

	case namespacesLoadedMsg:
		if msg.err != nil {
			m.state = menuError
			m.errMsg = msg.err.Error()
			return m, nil
		}
		m.namespaces = msg.namespaces
		m.state = menuNamespace
		m.cursor = 0
		return m, nil
	}

	return m, nil
}

var mainMenuItems = []string{"Start Game", "Select Namespace", "Exit"}

func (m MenuModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(mainMenuItems)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0: // Start Game
			m.k8sClient.SetNamespace(m.namespace)
			gameModel := NewGameModel(m.k8sClient, m.namespace, m.theme, m.width, m.height, m.kubeconfigPath)
			return gameModel, gameModel.Init()
		case 1: // Select Namespace
			return m, fetchNamespacesCmd(m.k8sClient)
		case 2: // Exit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) updateNamespace(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Item 0 is "all namespaces", then the actual namespaces follow
	totalItems := len(m.namespaces) + 1

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < totalItems-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor == 0 {
			m.namespace = ""
		} else {
			m.namespace = m.namespaces[m.cursor-1]
		}
		m.state = menuMain
		m.cursor = 0
	case "esc", "q":
		m.state = menuMain
		m.cursor = 0
	}
	return m, nil
}

func (m MenuModel) updateError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		m.state = menuConnecting
		return m, connectK8sCmd(m.kubeconfigPath)
	case "q", "esc":
		return m, tea.Quit
	}
	return m, nil
}

func (m MenuModel) View() string {
	if m.width == 0 {
		return ""
	}

	theme := m.theme

	titleArt := `
     _        _ _   _         _       _   ___
 ___| |_ __ _| | |_(_)_ _  __| |___  (_) |  _|
(_-<| ' / _  | / /| | ' \/ _  / / /  | | |_| |
/__/|_||_\__,_|_\_\|_|_||_\__,_\_\_\  |_| |___|
`

	title := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		Render(titleArt)

	subtitle := lipgloss.NewStyle().
		Foreground(theme.AccentSoft).
		Render("  snake your way through kubernetes pods")

	var body string

	switch m.state {
	case menuConnecting:
		body = lipgloss.NewStyle().
			Foreground(theme.Dim).
			Italic(true).
			Render(fmt.Sprintf("  connecting to cluster via %s ...", m.kubeconfigPath))

	case menuError:
		errBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Error).
			Foreground(theme.Error).
			Padding(1, 2).
			Render(fmt.Sprintf("Failed to connect:\n\n%s", m.errMsg))
		controls := lipgloss.NewStyle().
			Foreground(theme.Dim).
			Render("\n  [r] retry  [q] quit")
		body = errBox + controls

	case menuMain:
		body = m.viewMainMenu()

	case menuNamespace:
		body = m.viewNamespaceMenu()
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		subtitle,
		"",
		body,
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m MenuModel) viewMainMenu() string {
	theme := m.theme

	clusterInfo := lipgloss.NewStyle().
		Foreground(theme.Dim).
		Render(fmt.Sprintf("  cluster: %s", m.clusterName))

	nsLabel := "all"
	if m.namespace != "" {
		nsLabel = m.namespace
	}
	nsInfo := lipgloss.NewStyle().
		Foreground(theme.Dim).
		Render(fmt.Sprintf("  namespace: %s", nsLabel))

	configInfo := lipgloss.NewStyle().
		Foreground(theme.Dim).
		Render(fmt.Sprintf("  kubeconfig: %s", m.kubeconfigPath))

	var items []string
	for i, item := range mainMenuItems {
		if i == m.cursor {
			cursor := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("> ")
			label := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(item)
			items = append(items, "  "+cursor+label)
		} else {
			label := lipgloss.NewStyle().Foreground(theme.Foreground).Render("  "+item)
			items = append(items, "  "+label)
		}
	}

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Render(strings.Join(items, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left,
		clusterInfo,
		nsInfo,
		configInfo,
		"",
		menu,
		"",
		lipgloss.NewStyle().Foreground(theme.Dim).Render("  [j/k] navigate  [enter] select"),
	)
}

func (m MenuModel) viewNamespaceMenu() string {
	theme := m.theme

	header := lipgloss.NewStyle().
		Foreground(theme.AccentSoft).
		Bold(true).
		Render("  Select Namespace")

	var items []string

	// "All namespaces" option first
	allLabel := "all namespaces"
	if m.cursor == 0 {
		cursor := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("> ")
		label := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(allLabel)
		items = append(items, cursor+label)
	} else {
		items = append(items, "  "+lipgloss.NewStyle().Foreground(theme.Foreground).Render(allLabel))
	}

	for i, ns := range m.namespaces {
		idx := i + 1
		if idx == m.cursor {
			cursor := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("> ")
			label := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render(ns)
			items = append(items, cursor+label)
		} else {
			items = append(items, "  "+lipgloss.NewStyle().Foreground(theme.Foreground).Render(ns))
		}
	}

	// Show a scrollable window if there are many namespaces
	visible := items
	maxVisible := 15
	if len(items) > maxVisible {
		start := m.cursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end := start + maxVisible
		if end > len(items) {
			end = len(items)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
		visible = items[start:end]
	}

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Render(strings.Join(visible, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		menu,
		"",
		lipgloss.NewStyle().Foreground(theme.Dim).Render("  [j/k] navigate  [enter] select  [esc] back"),
	)
}

func connectK8sCmd(kubeconfigPath string) tea.Cmd {
	return func() tea.Msg {
		client, err := k8s.NewClient(kubeconfigPath, "")
		return k8sConnectedMsg{client: client, err: err}
	}
}

func fetchNamespacesCmd(client *k8s.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ns, err := client.ListNamespaces(ctx)
		return namespacesLoadedMsg{namespaces: ns, err: err}
	}
}
