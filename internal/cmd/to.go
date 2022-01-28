package cmd

import (
	"github.com/m-porter/jumper/internal/tui_v2"
	"github.com/spf13/cobra"
)

func ToCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to",
		Short: "Display projects in an intractable list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui_v2.Run()
		},
	}
}
