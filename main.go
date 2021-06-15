package main

import (
	"fmt"
	"github.com/m-porter/jumper/internal/jumper"
	"github.com/spf13/cobra"
	"os"
)

var flagAnalyze bool

func execute(cmd *cobra.Command, args []string) error {
	if flagAnalyze {
		jumper.Analyze()
		return nil
	}

	return jumper.Run(args)
}

func main() {
	cobra.OnInitialize(jumper.Init)

	jumperCmd := &cobra.Command{
		Use:   "jumper",
		Short: "Seamlessly jump between projects on your machine.",
		RunE:  execute,
	}

	jumperCmd.Flags().BoolVarP(&flagAnalyze, "analyze", "a", false, "Only run analyzer and exit")

	if err := jumperCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
