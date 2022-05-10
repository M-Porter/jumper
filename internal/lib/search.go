package lib

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	levenshtein "github.com/ka-weihe/fast-levenshtein"

	"go.uber.org/zap"

	"github.com/m-porter/jumper/internal/logger"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type match struct {
	value    string
	distance int
}

func normalizeString(in string) string {
	t := transform.Chain(norm.NFC, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, in)
	if err != nil {
		out = in
	}
	return out
}

func FuzzySearchSlice(search []string, term string) []string {
	term = normalizeString(term)

	re := regexp.MustCompile(
		fmt.Sprintf(
			"(?i).*%s.*",
			strings.Join(strings.Split(term, " "), ".*"),
		),
	)

	logger.Log("regexp", zap.Any("regexp", re.String()))

	var matches []match
	for _, s := range search {
		doc := normalizeString(filepath.Base(s))
		if re.MatchString(doc) {
			distance := levenshtein.Distance(doc, term)
			matches = append(matches, match{
				value:    s,
				distance: distance,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	var results []string
	for _, m := range matches {
		results = append(results, m.value)
	}

	return results
}
