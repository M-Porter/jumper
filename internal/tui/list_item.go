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

func (li *ListItem) LabelForStyle(style ListStyle) string {
	switch style {
	case ListStyleDetailed:
		return fmt.Sprintf("%s (%s)", li.Base, li.Dir)
	case ListStyleLong:
		return li.Path
	case ListStyleShort:
		fallthrough
	default:
		return li.Base
	}
}
