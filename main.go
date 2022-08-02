package main

import (
	"fmt"
	"os"
	"time"

	"github.com/m-porter/jumper/internal/cmd"
	"github.com/m-porter/jumper/internal/config"
	"github.com/spf13/cobra"
)

// set by goreleaser ldflags at build time
var version = "development"
var commit = ""
var date = time.Now().Format(time.RFC3339)
var builtBy = os.Getenv("USER")

func main() {
	cobra.OnInitialize(config.Init)

	rootCmd := cmd.RootCmd(cmd.RootCmdOptions{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	})

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
