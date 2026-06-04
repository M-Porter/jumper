package tui

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
)

func TestSearchBoxComponent(t *testing.T) {
	testCases := []struct {
		value string
		found int
		total int
	}{
		{"", 0, 0},
		{"", 0, 10},
		{"", 2, 10},
		{"", 10, 10},
		{"hello", 2, 10},
		{"some/path/to/file", 5, 42},
		{"foo", 1, 100},
		{"foo", 100, 100},
	}

	for _, tc := range testCases {
		width := 100

		searchBox := SearchBoxComponent(tc.value, tc.found, tc.total, width)
		plain := ansi.Strip(searchBox)

		assert.Equal(t, width, lipgloss.Width(searchBox))

		// value starts at the 3rd character (rune index 2), after "❯ "
		if tc.value != "" {
			assert.True(t, strings.HasPrefix(string([]rune(plain)[2:]), tc.value), "value should start at character index 2")
		}

		totalStr := fmt.Sprintf("/ %d", tc.total)
		assert.True(t, strings.Contains(plain, totalStr), "value should contain total")

		foundStr := fmt.Sprintf("%d /", tc.found)
		assert.True(t, strings.Contains(plain, foundStr), "value should contain found")
	}
}
