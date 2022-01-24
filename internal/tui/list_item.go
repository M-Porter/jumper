package tui

import (
	"fmt"
	"path/filepath"
)

func pathsToListItems(paths []string) []ListItem {
	var r []ListItem
	for _, path := range paths {
		r = append(r, ListItem{
			Path: path,
			Base: filepath.Base(path),
			Dir:  filepath.Dir(path),
		})
	}
	return r
}

type ListItem struct {
	Path string
	Base string
	Dir  string
}

func (li *ListItem) Format(style ListStyle, selected bool) string {
	var line string

	switch style {
	case ListStyleDetailed:
		line = fmt.Sprintf("%s (%s)", li.Base, li.Dir)
	case ListStyleLong:
		line = li.Path
	case ListStyleShort:
		fallthrough
	default:
		line = li.Base
	}

	spaceChar := " "
	if selected {
		spaceChar = ">"
		line = ColorBgGray.Sprintf(" %s ", line)
	} else {
		line = fmt.Sprintf(" %s", line)
	}

	spaceChar = ColorBgGray.Sprintf("%s", ColorFgRed.Sprint(spaceChar))

	line = ColorBgDefault.Sprintf("%s%s", spaceChar, line)

	return line
}
