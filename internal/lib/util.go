package lib

import (
	"regexp"
	"strings"
)

func QuoteParts(parts []string) []string {
	var escaped []string
	for _, part := range parts {
		escaped = append(escaped, regexp.QuoteMeta(part))
	}
	return escaped
}

func RegexpJoinPartsOr(parts []string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(QuoteParts(parts), "|"))
}

func RemoveDuplicates(dirs []string) []string {
	set := make(map[string]struct{})
	var r []string
	for _, dir := range dirs {
		if _, ok := set[dir]; !ok {
			r = append(r, dir)
			set[dir] = struct{}{}
		}
	}
	return r
}
