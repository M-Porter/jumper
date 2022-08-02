package main

import (
	"fmt"
	"os"

	"github.com/m-porter/jumper/internal/cmd"
	"github.com/m-porter/jumper/internal/config"
	"github.com/spf13/cobra"
)

// set by goreleaser ldflags at build time
//
// test locally with:
//    GORELEASER_CURRENT_TAG=v4.2.0 goreleaser build --single-target --snapshot --rm-dist
//    ./dist/jumper_{os}_{arch}_v1/jumper version
var version = "development"
var commit = "development"
var date = "development"
var builtBy = "development"

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
