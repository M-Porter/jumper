package main

import (
	"fmt"
	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
	"os"
)

func executeCmdTo(cmd *cobra.Command, args []string) error {
	return core.Run(args)
}

func executeCmdAnalyze(cmd *cobra.Command, args []string) {
	core.Analyze()
}

func main() {
	cobra.OnInitialize(core.Init)

	jumperCmd := &cobra.Command{
		Use:   "jumper",
		Short: "Seamlessly jump between projects on your machine.",
	}

	toCmd := &cobra.Command{
		Use:   "to",
		Short: "Run the jumper TUI.",
		RunE:  executeCmdTo,
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze and cache projects.",
		Run:   executeCmdAnalyze,
	}

	// todo: install cmd

	jumperCmd.AddCommand(toCmd, analyzeCmd)

	if err := jumperCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
