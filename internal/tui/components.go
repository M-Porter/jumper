package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/m-porter/jumper/internal/theme"
)

// InputComponent generates the string for the input line based on an input value.
func InputComponent(value string) string {
	pointerStyle := lipgloss.NewStyle().Foreground(theme.Mauve).Bold(true)
	return fmt.Sprintf("%s %s", pointerStyle.Render(theme.Pointer), value)
}

// SearchBoxComponent renders the search bar with input on the left and "found / total" on the right.
func SearchBoxComponent(value string, found int, total int, width int) string {
	input := " " + InputComponent(value)

	foundStyle := lipgloss.NewStyle().Foreground(theme.Green).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)

	stats := foundStyle.Render(fmt.Sprintf("%d", found)) + dimStyle.Render(fmt.Sprintf(" / %d ", total))

	// 2 bc of border
	gap := width - lipgloss.Width(input) - 2
	gap = max(0, gap)

	right := lipgloss.NewStyle().
		Width(gap).
		Align(lipgloss.Right).
		Render(stats)

	return lipgloss.NewStyle().BorderForeground(theme.Surface1).
		Border(lipgloss.RoundedBorder()).
		Render(input + right)
}

// ProjectRowComponent renders the individual row for each project
func ProjectRowComponent(value string, selected bool, truncatePath bool) string {
	pathStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)

	var selectedStyle lipgloss.Style
	var projectStyle lipgloss.Style
	selectedPointer := " "
	if selected {
		selectedStyle = lipgloss.NewStyle().Foreground(theme.Blue).Bold(true)
		projectStyle = lipgloss.NewStyle().Foreground(theme.Blue).Bold(true)
		selectedPointer = fmt.Sprintf("%s", theme.Pointer)
	} else {
		selectedStyle = lipgloss.NewStyle().Foreground(theme.Surface2)
		projectStyle = lipgloss.NewStyle().Foreground(theme.Text)
	}

	dir := filepath.Dir(value)
	proj := filepath.Base(value)

	if homedir, err := os.UserHomeDir(); err == nil {
		dir = strings.Replace(dir, homedir, "~", 1)
	}

	var dirStr string
	pathParts := strings.Split(dir, string(os.PathSeparator))
	if len(pathParts) > 2 && truncatePath {
		dirStr = strings.Join(pathParts[:2], string(os.PathSeparator)) + "/" + theme.Ellipsis + "/"
	} else {
		dirStr = dir + "/"
	}

	// the following bit fixes weird looking spacing when the folder glyph isn't
	// included basically changes "❯   ~/foo" to "❯ ~/foo"
	folderGlyph := strings.TrimLeft(fmt.Sprintf("%s  ", theme.Folder()), " ")

	return fmt.Sprintf("%s %s%s%s", selectedStyle.Render(selectedPointer), selectedStyle.Render(folderGlyph), pathStyle.Render(dirStr), projectStyle.Render(proj))
}

type StatusBarParams struct {
	LeftContents  []string
	RightContents []string
}

// StatusBarComponent renders the bottom status bar that contains things like keyboard shortcuts and helpful hint
// messages
func StatusBarComponent(params StatusBarParams, width int) string {
	leftContents := strings.Join(params.LeftContents, "  ")
	rightContents := strings.Join(params.RightContents, "  ")

	gap := width - lipgloss.Width(leftContents) - 2
	gap = max(0, gap)

	rightContents = lipgloss.NewStyle().
		Width(gap).
		Align(lipgloss.Right).
		Render(rightContents)

	return lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(1).
		BorderForeground(theme.Surface1).
		Border(lipgloss.RoundedBorder(), true, false, false, false).
		Render(leftContents + rightContents)
}

func KeyHelpComponent(keys, hint string) string {
	keysStyle := lipgloss.NewStyle().Foreground(theme.Blue)
	hintStyle := lipgloss.NewStyle().Foreground(theme.Overlay1)

	var rendered []string
	if keys != "" {
		rendered = append(rendered, keysStyle.Render(keys))
	}
	if hint != "" {
		rendered = append(rendered, hintStyle.Render(hint))
	}

	return strings.Join(rendered, " ")
}
