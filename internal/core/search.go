package core

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
)

func filterDirectories(data []string, term string) []string {
	matches := fuzzy.RankFindNormalizedFold(term, data)
	sort.Stable(matches)
	var results []string
	for _, match := range matches {
		results = append(results, match.Target)
	}
	return results
}
