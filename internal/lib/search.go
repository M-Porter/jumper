package lib

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
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

// FuzzySearchSlice does a fuzzy search on haystack of paths given a needle. Searching is only done against the
// base name of a given path. So given "/usr/dev/my-project", and "/usr/dev/dev-tools", a needle of "dev"
// will only yield "/usr/dev/dev-tools".
//
// If the given needle value yield an invalid regex, the entire haystack is returned.
//
// The returned slice will be sorted based on Levenshtein distance.
func FuzzySearchSlice(haystack []string, needle string) []string {
	searchResults := make(chan match, len(haystack))
	var wg sync.WaitGroup

	needle = normalizeString(needle)

	re, err := regexp.Compile(
		fmt.Sprintf(
			"(?i).*%s.*",
			strings.Join(strings.Split(needle, " "), ".*"),
		),
	)
	if err != nil {
		logger.Log("Failed to compile regex", zap.String("regex", needle), zap.Error(err))
		return haystack
	}

	logger.Log("regexp", zap.Any("regexp", re.String()))

	for _, s := range haystack {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()

			doc := normalizeString(filepath.Base(value))
			if re.MatchString(doc) {
				distance := levenshtein.Distance(doc, needle)
				searchResults <- match{
					value:    value,
					distance: distance,
				}
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(searchResults)
	}()

	// collect
	matches := make([]match, 0, len(haystack))
	for sr := range searchResults {
		matches = append(matches, sr)
	}

	// sort
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	// transform
	results := make([]string, len(matches))
	for i, m := range matches {
		results[i] = m.value
	}

	return results
}
