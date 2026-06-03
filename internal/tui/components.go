package tui

import (
	"fmt"

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

	foundStyle := lipgloss.NewStyle().Foreground(theme.Green).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)

	stats := foundStyle.Render(fmt.Sprintf("%d", found)) + dimStyle.Render(fmt.Sprintf(" / %d ", total))

	gap := width - lipgloss.Width(input)
	gap = max(0, gap)

	right := lipgloss.NewStyle().
		Width(gap).
		Align(lipgloss.Right).
		Render(stats)

	return input + right
}
