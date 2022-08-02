package lib

import (
	"errors"
	"os"
	"os/exec"
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

func AbsValue[V float64](val V) V {
	if val < 0 {
		return -val
	}
	return val
}

func OpenEditor(filePath string) error {
	var editorCmd *exec.Cmd

	if editor := os.Getenv("EDITOR"); editor != "" {
		editorCmd = exec.Command(editor, filePath)
	} else {
		editorCmd = exec.Command("which", "vim", "nano")
		editorCmd.Stdin = os.Stdin
		out, err := editorCmd.Output()
		if err != nil {
			return errors.New("could not determine which editor to use")
		}
		e := strings.Split(string(out), "\n")[0]
		editorCmd = exec.Command(e, filePath)
	}

	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	return editorCmd.Run()
}
