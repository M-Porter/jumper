package cmd

import "github.com/spf13/cobra"

func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install jumper helpers to your environment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo this
			return nil
		},
	}
}
