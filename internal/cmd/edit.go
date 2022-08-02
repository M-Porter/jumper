package cmd

import (
	"github.com/m-porter/jumper/internal/config"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/spf13/cobra"
)

func EditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit the jumper config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.OpenEditor(config.Filepath())
		},
	}
}
