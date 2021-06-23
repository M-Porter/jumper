package main

import (
	"fmt"
	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
	"os"
)

// jumper
var rootCmd = &cobra.Command{
	Use:   "jumper",
	Short: "Seamlessly jump between projects on your machine.",
}

// jumper to
var toCmd = &cobra.Command{
	Use:   "to",
	Short: "Run the jumper TUI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Run(args)
	},
}

// jumper analyze
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze and cache projects.",
	Run: func(cmd *cobra.Command, args []string) {
		core.Analyze()
	},
}

func main() {
	cobra.OnInitialize(core.Init)

	// todo: install cmd
	rootCmd.AddCommand(toCmd, analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
