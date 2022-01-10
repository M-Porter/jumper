package cmd

import "github.com/spf13/cobra"

func SetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Outputs setup instructions on how to most effectively use jumper on your system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("todo")
		},
	}
}
