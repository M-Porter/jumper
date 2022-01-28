package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const lineIndicator = "❯"

const (
	colorZinc700 = lipgloss.Color("#3F3F46")
	colorZinc500 = lipgloss.Color("#71717A")
	colorRed     = lipgloss.Color("#EF4444")
	colorBlue    = lipgloss.Color("#0EA5E9")
)

var (
	bgGrayStyle          = lipgloss.NewStyle().Background(colorZinc700).Bold(true)
	indicatorStyle       = lipgloss.NewStyle().Bold(true).Foreground(colorRed).Background(colorZinc700)
	detailDimStyleBgGray = lipgloss.NewStyle().Bold(true).Foreground(colorZinc500).Background(colorZinc700)
	detailDimStyle       = lipgloss.NewStyle().Foreground(colorZinc500)
	inputArrowStyle      = lipgloss.NewStyle().Bold(true).Foreground(colorBlue)

	selectedRowIndicatorPart = indicatorStyle.Render(lineIndicator)
	inputIndicatorPart       = inputArrowStyle.Render(lineIndicator)
	beamPart                 = bgGrayStyle.Render(" ")
)

type listStyle int

const (
	listStyleShort listStyle = iota
	listStyleLong
	listStyleDetailed
)

var (
	listStyles = []listStyle{listStyleShort, listStyleLong, listStyleDetailed}
)

func (ls listStyle) format(item listItem, selected bool) string {
	switch ls {
	case listStyleLong:
		return formatListStyleLong(item, selected)
	case listStyleDetailed:
		return formatListStyleDetailed(item, selected)
	case listStyleShort:
		fallthrough
	default:
		return formatListStyleShort(item, selected)
	}
}

func formatListStyleShort(item listItem, selected bool) string {
	var line string

	if selected {
		infoPart := bgGrayStyle.Render(fmt.Sprintf(" %s ", item.Base))
		line = fmt.Sprintf("%s%s", selectedRowIndicatorPart, infoPart)
	} else {
		infoPart := fmt.Sprintf(" %s ", item.Base)
		line = fmt.Sprintf("%s%s", beamPart, infoPart)
	}

	return line
}

func formatListStyleLong(item listItem, selected bool) string {
	var line string

	if selected {
		longPart := detailDimStyleBgGray.Render(fmt.Sprintf(" %s/", item.Dir))
		shortPart := bgGrayStyle.Render(item.Base)
		line = fmt.Sprintf("%s%s%s", selectedRowIndicatorPart, longPart, shortPart)
	} else {
		longPart := detailDimStyle.Render(fmt.Sprintf(" %s/", item.Dir))
		shortPart := item.Base
		line = fmt.Sprintf("%s%s%s", selectedRowIndicatorPart, longPart, shortPart)
		line = fmt.Sprintf("%s%s%s", beamPart, longPart, shortPart)
	}

	return line
}

func formatListStyleDetailed(item listItem, selected bool) string {
	var line string

	if selected {
		detailPart := detailDimStyleBgGray.Render(fmt.Sprintf("(%s) ", item.Dir))
		infoPart := bgGrayStyle.Render(fmt.Sprintf(" %s %s", item.Base, detailPart))
		line = fmt.Sprintf("%s%s", selectedRowIndicatorPart, infoPart)
	} else {
		detailPart := detailDimStyle.Render(fmt.Sprintf("(%s) ", item.Dir))
		infoPart := fmt.Sprintf(" %s %s", item.Base, detailPart)
		line = fmt.Sprintf("%s%s", beamPart, infoPart)
	}

	return line
}
