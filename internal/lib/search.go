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

func FuzzySearchSlice(search []string, term string) []string {
	searchResults := make(chan match, len(search))
	var wg sync.WaitGroup

	term = normalizeString(term)

	re := regexp.MustCompile(
		fmt.Sprintf(
			"(?i).*%s.*",
			strings.Join(strings.Split(term, " "), ".*"),
		),
	)

	logger.Log("regexp", zap.Any("regexp", re.String()))

	for _, s := range search {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()

			doc := normalizeString(filepath.Base(value))
			if re.MatchString(doc) {
				distance := levenshtein.Distance(doc, term)
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
	matches := make([]match, 0, len(search))
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
