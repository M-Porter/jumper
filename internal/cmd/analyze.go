package cmd

import (
	"fmt"

	"github.com/m-porter/jumper/internal/core"
	"github.com/spf13/cobra"
)

func AnalyzeCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "analyze",
		Aliases: []string{"setup"},
		Short:   "Search for and cache projects",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Analyzing...")
			app := core.NewApp()
			app.Analyze()
			dirs := app.Directories

			fmt.Printf("Projects found: %d\n", len(dirs))
			for _, dir := range dirs {
				fmt.Printf("  - %s\n", dir)
			}
		},
	}
}
