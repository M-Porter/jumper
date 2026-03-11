package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	selectedRowIndicatorPart string
	inputIndicatorPart       string
)

func initIndicators(indicator string) {
	selectedRowIndicatorPart = indicatorStyle.Render(indicator)
	inputIndicatorPart = inputArrowStyle.Render(indicator)
}

var (
	// https://catppuccin.com/palette/
	colorSurface0 = lipgloss.AdaptiveColor{
		Light: "#ccd0da",
		Dark:  "#313244",
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
	colorOverlay2 = lipgloss.AdaptiveColor{
		Light: "#7c7f93",
		Dark:  "#9399b2",
	}

	bgGrayStyle          = lipgloss.NewStyle().Background(colorSurface0).Bold(true)
	indicatorStyle       = lipgloss.NewStyle().Bold(true).Foreground(colorMaroon).Background(colorSurface0)
	detailDimStyleBgGray = lipgloss.NewStyle().Bold(true).Foreground(colorOverlay2).Background(colorSurface0)
	detailDimStyle       = lipgloss.NewStyle().Foreground(colorOverlay2)
	inputArrowStyle      = lipgloss.NewStyle().Bold(true).Foreground(colorLavender)
	textStyle            = lipgloss.NewStyle().Foreground(colorText)

	beamPart = bgGrayStyle.Render(" ")
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
