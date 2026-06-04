package cmd

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/m-porter/jumper/internal/core"
	"github.com/m-porter/jumper/internal/theme"
	"github.com/m-porter/jumper/internal/tui"
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

			fmt.Printf("Projects found: %s\n", lipgloss.NewStyle().Foreground(theme.Green).Bold(true).Render(fmt.Sprintf("%d", len(dirs))))
			for _, dir := range dirs {
				fmt.Println(tui.ProjectRowComponent(dir, false, false))
			}
		},
	}
}
