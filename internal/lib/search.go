package lib

import (
	"sort"

	"github.com/m-porter/jumper/internal/logger"
	"go.uber.org/zap"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func FuzzySearchSlice(data []string, term string) []string {
	matches := fuzzy.RankFindNormalizedFold(term, data)
	sort.Sort(matches)
	logger.Log("FuzzySearch", zap.Any("input_data", data), zap.Any("matches", matches))
	var results []string
	for _, match := range matches {
		results = append(results, match.Target)
	}
	return results
}
