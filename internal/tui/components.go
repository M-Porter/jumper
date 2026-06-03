package tui

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
	"github.com/m-porter/jumper/internal/theme"
)

// InputComponent generates the string for the input line based on an input value.
func InputComponent(value string) string {
	pointerStyle := lipgloss.NewStyle().Foreground(theme.Mauve).Bold(true)
	return fmt.Sprintf("%s %s", pointerStyle.Render(theme.Pointer), value)
}

// SearchBox renders the entire top bar of the TUI with the Input on the left and the "x / y" value on the right.
func SearchBox(value string, found int, total int, width int) string {
	fb := flexbox.NewHorizontal(width, 1)

	foundStyle := lipgloss.NewStyle().Foreground(theme.Green).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)
	stats := fmt.Sprintf(
		"%s%s",
		foundStyle.Render(fmt.Sprintf("%d", found)),
		dimStyle.Render(fmt.Sprintf(" / %d", total)),
	)

	inputCell := flexbox.NewCell(1, 1).SetContent(InputComponent(value))
	statsCell := flexbox.NewCell(1, 1).
		SetStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Right)).
		SetContent(stats)

	cols := []*flexbox.Column{
		fb.NewColumn().AddCells(inputCell),
		fb.NewColumn().AddCells(statsCell),
	}

	fb.AddColumns(cols)

	return fb.Render()
}
