package main

import (
	"fmt"
	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
	"os"
)

var flagAnalyze bool

func execute(cmd *cobra.Command, args []string) error {
	if flagAnalyze {
		core.Analyze()
		return nil
	}
	return core.Run(args)
}

func main() {
	cobra.OnInitialize(core.Init)

	jumperCmd := &cobra.Command{
		Use:   "core",
		Short: "Seamlessly jump between projects on your machine.",
		RunE:  execute,
	}

	jumperCmd.Flags().BoolVarP(&flagAnalyze, "analyze", "a", false, "Only run analyzer and exit")

	if err := jumperCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
