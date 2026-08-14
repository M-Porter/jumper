package theme

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/m-porter/jumper/internal/config"
)

const (
	Pointer    = "❯"
	Enter      = "↵"
	Ellipsis   = "…"
	ArrowRight = "→"
	ArrowUp    = "↑"
	ArrowDown  = "↓"
	folder     = "\ue5ff"
	Checkmark  = "✓"
	Bullet     = "·"
)

type spinner struct {
	current int
	frames  []string
}

type SpinnerIter interface {
	Next() string
}

func Folder() string {
	if config.Get().NoNerdFont {
		return ""
	}
	return folder
}

func NewSpinner() SpinnerIter {
	return &spinner{
		current: 0,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (s *spinner) Next() string {
	f := s.frames[s.current]
	next := s.current + 1
	if next >= len(s.frames) {
		s.current = 0
	} else {
		s.current = next
	}
	return f
}

// https://catppuccin.com/palette/
var (
	Blue     = compat.AdaptiveColor{Light: lipgloss.Color("#1e66f5"), Dark: lipgloss.Color("#89b4fa")}
	Mauve    = compat.AdaptiveColor{Light: lipgloss.Color("#8839ef"), Dark: lipgloss.Color("#cba6f7")}
	Yellow   = compat.AdaptiveColor{Light: lipgloss.Color("#df8e1d"), Dark: lipgloss.Color("#f9e2af")}
	Sky      = compat.AdaptiveColor{Light: lipgloss.Color("#209fb5"), Dark: lipgloss.Color("#89dceb")}
	Green    = compat.AdaptiveColor{Light: lipgloss.Color("#40a02b"), Dark: lipgloss.Color("#a6e3a1")}
	Text     = compat.AdaptiveColor{Light: lipgloss.Color("#4c4f69"), Dark: lipgloss.Color("#cdd6f4")}
	Overlay1 = compat.AdaptiveColor{Light: lipgloss.Color("#8c8fa1"), Dark: lipgloss.Color("#7f849c")}
	Surface0 = compat.AdaptiveColor{Light: lipgloss.Color("#ccd0da"), Dark: lipgloss.Color("#313244")}
	Surface1 = compat.AdaptiveColor{Light: lipgloss.Color("#bcc0cc"), Dark: lipgloss.Color("#45475a")}
	Surface2 = compat.AdaptiveColor{Light: lipgloss.Color("#acb0be"), Dark: lipgloss.Color("#585b70")}
	Mantle   = compat.AdaptiveColor{Light: lipgloss.Color("#e6e9ef"), Dark: lipgloss.Color("#181825")}
	Base     = compat.AdaptiveColor{Light: lipgloss.Color("#eff1f5"), Dark: lipgloss.Color("#1e1e2e")}
)
