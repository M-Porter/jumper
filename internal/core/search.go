package core

import "github.com/sahilm/fuzzy"

func filterDirectories(data []string, term string) fuzzy.Matches {
	return fuzzy.Find(term, data)
}
