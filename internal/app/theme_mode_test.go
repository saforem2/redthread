package app

import "testing"

func TestParseThemeMode(t *testing.T) {
	for _, tc := range []struct {
		in      string
		want    ThemeMode
		wantErr bool
	}{
		{"auto", ThemeAuto, false},
		{"", ThemeAuto, false},
		{"AUTO", ThemeAuto, false},
		{"  auto  ", ThemeAuto, false},
		{"light", ThemeLight, false},
		{"Light", ThemeLight, false},
		{"dark", ThemeDark, false},
		{"DARK", ThemeDark, false},
		{"blue", ThemeAuto, true},
		{"lite", ThemeAuto, true},
	} {
		got, err := ParseThemeMode(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseThemeMode(%q) = %v, nil; want an error", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseThemeMode(%q) returned error %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseThemeMode(%q) = %v; want %v", tc.in, got, tc.want)
		}
	}
}

func TestThemeModeString(t *testing.T) {
	for mode, want := range map[ThemeMode]string{
		ThemeAuto:  "auto",
		ThemeDark:  "dark",
		ThemeLight: "light",
	} {
		if got := mode.String(); got != want {
			t.Errorf("ThemeMode(%d).String() = %q; want %q", mode, got, want)
		}
	}
}

// An explicit mode must never consult the terminal — that is the whole
// point of the override, since the OSC 11 probe is what fails under tmux
// and over ssh.
func TestResolveThemeHonorsExplicitMode(t *testing.T) {
	probed := false
	probe := func() bool {
		probed = true
		return true // claim "dark" regardless
	}

	if got := resolveTheme(ThemeLight, probe); !got {
		t.Error("resolveTheme(ThemeLight) = false; want true (light)")
	}
	if probed {
		t.Error("resolveTheme(ThemeLight) probed the terminal; an explicit mode must not")
	}

	probed = false
	if got := resolveTheme(ThemeDark, probe); got {
		t.Error("resolveTheme(ThemeDark) = true; want false (dark)")
	}
	if probed {
		t.Error("resolveTheme(ThemeDark) probed the terminal; an explicit mode must not")
	}
}

func TestResolveThemeAutoUsesProbe(t *testing.T) {
	if got := resolveTheme(ThemeAuto, func() bool { return true }); got {
		t.Error("resolveTheme(ThemeAuto) with a dark terminal = true; want false (dark)")
	}
	if got := resolveTheme(ThemeAuto, func() bool { return false }); !got {
		t.Error("resolveTheme(ThemeAuto) with a light terminal = false; want true (light)")
	}
}

// When the probe cannot reach a terminal it reports a zero-value (black)
// background, i.e. dark. That is the pre-existing behavior and the safe
// fallback: an unreadable-but-familiar dark board beats guessing light.
func TestResolveThemeAutoFallsBackToDark(t *testing.T) {
	if got := resolveTheme(ThemeAuto, hasDarkBackground); got {
		// hasDarkBackground on a non-TTY (as under `go test`) reports dark.
		t.Error("resolveTheme(ThemeAuto) on a non-TTY = true; want false (dark fallback)")
	}
}

func TestThemeModeFromEnv(t *testing.T) {
	for _, tc := range []struct {
		env  string
		flag string
		want ThemeMode
	}{
		{"", "", ThemeAuto},
		{"light", "", ThemeLight},
		{"dark", "", ThemeDark},
		{"garbage", "", ThemeAuto},    // invalid env is ignored, not fatal
		{"dark", "light", ThemeLight}, // an explicit flag beats the env
		{"light", "dark", ThemeDark},
		{"", "light", ThemeLight},
	} {
		got := themeModeFrom(tc.flag, tc.env)
		if got != tc.want {
			t.Errorf("themeModeFrom(flag=%q, env=%q) = %v; want %v",
				tc.flag, tc.env, got, tc.want)
		}
	}
}
