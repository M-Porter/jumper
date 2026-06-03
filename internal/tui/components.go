package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/m-porter/jumper/internal/theme"
)

// InputComponent generates the string for the input line based on an input value.
func InputComponent(value string) string {
	pointerStyle := lipgloss.NewStyle().Foreground(theme.Mauve).Bold(true)
	return fmt.Sprintf("%s %s", pointerStyle.Render(theme.Pointer), value)
}

// SearchBoxComponent renders the search bar with input on the left and "found / total" on the right.
func SearchBoxComponent(value string, found int, total int, width int) string {
	input := InputComponent(value)

	foundStr := fmt.Sprintf("%d", found)
	totalStr := fmt.Sprintf(" / %d", total)

	foundStyle := lipgloss.NewStyle().Foreground(theme.Green).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)

	stats := foundStyle.Render(foundStr) + dimStyle.Render(totalStr)

	gap := width - lipgloss.Width(input) - lipgloss.Width(stats)
	if gap < 0 {
		gap = 0
	}

	return input + strings.Repeat(" ", gap) + stats
}
