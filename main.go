package main

import (
	"fmt"
	"os"

	"github.com/m-porter/jumper/internal/cmd"
	"github.com/m-porter/jumper/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	cobra.OnInitialize(config.Init)

	if err := cmd.RootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
