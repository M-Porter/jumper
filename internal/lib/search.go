package lib

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"go.uber.org/zap"

	"github.com/m-porter/jumper/internal/logger"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalizeString(in string) string {
	t := transform.Chain(norm.NFC, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, in)
	if err != nil {
		out = in
	}
	return out
}

func FuzzySearchSlice(search []string, term string) []string {
	var matches []string

	term = normalizeString(term)

	re := regexp.MustCompile(
		fmt.Sprintf(
			"(?i).*%s.*",
			strings.Join(strings.Split(term, " "), ".*"),
		),
	)

	logger.Log("regexp", zap.Any("regexp", re.String()))

	for _, s := range search {
		doc := normalizeString(filepath.Base(s))
		if re.MatchString(doc) {
			matches = append(matches, s)
		}
	}

	return matches
}
