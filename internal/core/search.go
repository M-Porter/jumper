package core

import "github.com/sahilm/fuzzy"

func filterDirectories(dirs []Dir, term string) fuzzy.Matches {
	var data []string
	for _, d := range dirs {
		data = append(data, d.Label)
	}
	return fuzzy.Find(term, data)
}
