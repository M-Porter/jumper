package core

import "github.com/sahilm/fuzzy"

func filterDirectories(data []string, term string) []string {
	matches := fuzzy.Find(term, data)
	var results []string
	for _, match := range matches {
		results = append(results, match.Str)
	}
	return results
}
