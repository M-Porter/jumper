package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzySearchSlice(t *testing.T) {
	testCases := []struct {
		haystack []string
		needle   string
		expected []string
	}{
		{
			haystack: []string{"/some/directory/project", "/some/foo/dir", "/users/apples/oranges"},
			needle:   "dir",
			expected: []string{"/some/foo/dir"},
		},
		{
			// empty haystack
			haystack: []string{},
			needle:   "foo",
			expected: []string{},
		},
		{
			// needle matches nothing
			haystack: []string{"/a/foo", "/b/bar"},
			needle:   "xyz",
			expected: []string{},
		},
		{
			// match is case-insensitive
			haystack: []string{"/a/Foobar", "/b/baz"},
			needle:   "fOo",
			expected: []string{"/a/Foobar"},
		},
		{
			// only the basename is matched, not parent path segments
			haystack: []string{"/needle/project", "/other/needle"},
			needle:   "needle",
			expected: []string{"/other/needle"},
		},
		{
			// results are sorted by levenshtein distance, closest first
			haystack: []string{"/a/foobar", "/b/foo"},
			needle:   "foo",
			expected: []string{"/b/foo", "/a/foobar"},
		},
	}

	for _, tc := range testCases {
		results := FuzzySearchSlice(tc.haystack, tc.needle)
		assert.Equal(t, tc.expected, results)
	}
}
