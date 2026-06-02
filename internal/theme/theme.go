package theme

import "github.com/charmbracelet/lipgloss"

type Glyph string

const (
	Pointer    Glyph = "❯"
	Enter            = "↵"
	Ellipsis         = "…"
	ArrowRight       = "→"
	ArrowUp          = "↑"
	ArrowDown        = "↓"
	Folder           = "\ue5ff"
	Checkmark        = "✓"
	Bullet           = "·"
)

type spinner struct {
	current int
	frames  []string
}

type SpinnerIter interface {
	Next() string
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
	Blue     = lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#89b4fa"}
	Mauve    = lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#cba6f7"}
	Yellow   = lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f9e2af"}
	Sky      = lipgloss.AdaptiveColor{Light: "#209fb5", Dark: "#89dceb"}
	Green    = lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6e3a1"}
	Text     = lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cdd6f4"}
	Overlay1 = lipgloss.AdaptiveColor{Light: "#8c8fa1", Dark: "#7f849c"}
	Surface0 = lipgloss.AdaptiveColor{Light: "#ccd0da", Dark: "#313244"}
	Surface1 = lipgloss.AdaptiveColor{Light: "#bcc0cc", Dark: "#45475a"}
	Surface2 = lipgloss.AdaptiveColor{Light: "#acb0be", Dark: "#585b70"}
	Mantle   = lipgloss.AdaptiveColor{Light: "#e6e9ef", Dark: "#181825"}
	Base     = lipgloss.AdaptiveColor{Light: "#eff1f5", Dark: "#1e1e2e"}
)
