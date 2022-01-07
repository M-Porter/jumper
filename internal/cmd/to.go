package cmd

import (
	"github.com/m-porter/jumper/internal/core"
	"github.com/m-porter/jumper/internal/tui"
	"github.com/spf13/cobra"
)

func ToCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to",
		Short: "Run the jumper TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.New(core.NewApp(runInDebugMode)).Run()
		},
	}
}
