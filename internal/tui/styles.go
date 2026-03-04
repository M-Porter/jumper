package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	lineIndicator = "❯"

	// https://catppuccin.com/palette/
	colorSurface1 = lipgloss.AdaptiveColor{
		Light: "#bcc0cc",
		Dark:  "#45475a",
	}
	colorMaroon = lipgloss.AdaptiveColor{
		Light: "#e64553",
		Dark:  "#eba0ac",
	}
	colorLavender = lipgloss.AdaptiveColor{
		Light: "#7287fd",
		Dark:  "#b4befe",
	}
	colorText = lipgloss.AdaptiveColor{
		Light: "#4c4f69",
		Dark:  "#cdd6f4",
	}
	colorSubtext0 = lipgloss.AdaptiveColor{
		Light: "#6c6f85",
		Dark:  "#a6adc8",
	}

	bgGrayStyle          = lipgloss.NewStyle().Background(colorSurface1).Bold(true)
	indicatorStyle       = lipgloss.NewStyle().Bold(true).Foreground(colorMaroon).Background(colorSurface1)
	detailDimStyleBgGray = lipgloss.NewStyle().Bold(true).Foreground(colorSubtext0).Background(colorSurface1)
	detailDimStyle       = lipgloss.NewStyle().Foreground(colorSubtext0)
	inputArrowStyle      = lipgloss.NewStyle().Bold(true).Foreground(colorLavender)
	textStyle            = lipgloss.NewStyle().Foreground(colorText)

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

	baseInfoPart := textStyle.Render(fmt.Sprintf(" %s ", item.Base))
	if selected {
		line = fmt.Sprintf("%s%s", selectedRowIndicatorPart, bgGrayStyle.Render(baseInfoPart))
	} else {
		line = fmt.Sprintf("%s%s", beamPart, baseInfoPart)
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
