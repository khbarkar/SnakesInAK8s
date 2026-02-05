package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kristinb/snakeinak8/internal/k8s"
	"github.com/kristinb/snakeinak8/internal/ui"
)

func main() {
	kubeconfigFlag := flag.String("kubeconfig", "", "path to kubeconfig file (defaults to KUBECONFIG env or ~/.kube/config)")
	flag.Parse()

	kubeconfigPath := k8s.ResolveKubeconfig(*kubeconfigFlag)

	m := ui.NewMenuModel(kubeconfigPath)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
