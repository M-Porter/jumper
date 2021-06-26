package main

import (
	"errors"
	"fmt"
	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
	"os"
)

var (
	runInDebugMode bool

	resetCache  bool
	resetConfig bool
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
		return core.Run(runInDebugMode)
	},
}

// jumper analyze
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze and cache projects.",
	Run: func(cmd *cobra.Command, args []string) {
		core.RunAnalyzer(runInDebugMode)
	},
}

// jumper clear
var resetCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the jumper cache and/or config.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !resetCache && !resetConfig {
			return errors.New("--cache and/or --config flags required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// todo this
	},
}

// jumper install
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install jumper helpers to your environment.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// todo this
		return nil
	},
}

func main() {
	cobra.OnInitialize(core.Init)

	rootCmd.PersistentFlags().BoolVar(&runInDebugMode, "debug", false, "Run jumper in debug mode.")

	resetCmd.Flags().BoolVar(&resetCache, "cache", false, "Reset the jumper cache.")
	resetCmd.Flags().BoolVar(&resetConfig, "config", false, "Reset the jumper config.")

	rootCmd.AddCommand(
		toCmd,
		analyzeCmd,
		resetCmd,
		installCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
