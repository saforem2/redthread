package app

// app.go — top-level entry. Loads workspace (or seeds a fresh one) and
// launches the Bubble Tea program.

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// envTruthy reads the usual affirmative spellings for a boolean env var.
func envTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

// Run is the package's main entry point. cmd/redthread/main.go calls it.
func Run() {
	var fresh bool
	var vimFlag bool
	flag.BoolVar(&fresh, "fresh", false, "start with a fresh seeded workspace, ignore saved notes")
	flag.BoolVar(&vimFlag, "vim", false, "modal (vim) editing in the note card. Also settable with RT_VIM=1, and remembered once used.")
	flag.Parse()

	var ws *Workspace
	if !fresh {
		w, err := LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load notes (%v) — starting fresh\n", err)
		}
		ws = w
	}
	if ws == nil || len(ws.Boards) == 0 {
		ws = seedWorkspace()
	}
	// Modal editing: the flag or RT_VIM turns it on for this run and is
	// remembered; otherwise the saved preference stands.
	if vimFlag || envTruthy(os.Getenv("RT_VIM")) {
		ws.Vim = true
	}

	if active := ws.ActiveBoard(); active != nil {
		active.ApplyGlobalBorder()
	}

	m := initialModel(ws)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = SaveWorkspace(ws)
}
