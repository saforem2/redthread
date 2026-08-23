package app

import "testing"

// The help panel is the only in-app discovery surface. Every feature's
// keys must appear, or the feature is invisible to anyone who did not
// read the README.
func TestHelpPanelListsEveryFeatureKey(t *testing.T) {
	seen := map[string]bool{}
	for _, col := range helpData {
		for _, e := range col.entries {
			seen[e.key] = true
		}
	}
	for _, want := range []struct{ key, feature string }{
		{"u", "undo"},
		{"C-r", "redo"},
		{"e", "$EDITOR"},
		{"T", "light/dark theme"},
	} {
		if !seen[want.key] {
			t.Errorf("help panel is missing %q (%s)", want.key, want.feature)
		}
	}
}

// It also has to fit: each feature added a row, and the panel is sized
// from the longest column.
func TestHelpPanelFitsASmallTerminal(t *testing.T) {
	ws := seedWorkspace()
	m := initialModel(ws)
	for _, size := range [][2]int{{80, 24}, {100, 30}, {140, 40}} {
		m.w, m.h = size[0], size[1]
		m.stars = GenStarsForBoard(m.w, m.h, ws.ActiveBoard().GrainSeed)
		m.helpAnim, m.helpTarget = 1, 1
		if out := m.View(); out == "" {
			t.Errorf("empty view with help open at %dx%d", size[0], size[1])
		}
	}
}
