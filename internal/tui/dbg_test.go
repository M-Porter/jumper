package tui

import (
	"fmt"
	"testing"
	"github.com/charmbracelet/x/ansi"
)

func TestDebugPlain(t *testing.T) {
	s := SearchBoxComponent("hello", 2, 10, 100)
	plain := ansi.Strip(s)
	fmt.Printf("plain: %q\n", plain)
	fmt.Printf("last 20 chars: %q\n", plain[len(plain)-20:])
}
