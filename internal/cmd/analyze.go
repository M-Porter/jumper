package cmd

import (
	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
)

func AnalyzeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze",
		Short: "search for and cache projects",
		Run: func(cmd *cobra.Command, args []string) {
			app := core.NewApp()
			app.Setup()
			app.Analyze()
		},
	}
}
