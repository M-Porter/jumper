package main

import (
	"fmt"
	"github.com/m-porter/jumper/internal/cmd"
	"github.com/m-porter/jumper/internal/config"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cobra.OnInitialize(config.Init)

	if err := cmd.RootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
