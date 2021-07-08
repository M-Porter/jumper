package cmd

import (
	"github.com/spf13/cobra"
)

func ClearCmd() *cobra.Command {
	clearCacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Clear the jumper cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	clearConfigCmd := &cobra.Command{
		Use:   "cache",
		Short: "Clear the jumper config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	clearAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Clear the jumper config and cache.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the jumper config or cache.",
	}

	cmd.AddCommand(
		clearAllCmd,
		clearCacheCmd,
		clearConfigCmd,
	)

	return cmd
}
