package cmd

import (
	"github.com/spf13/cobra"
)

func ClearCmd() *cobra.Command {
	clearCacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Clear the jumper cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("todo")
		},
	}

	clearConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Clear the jumper config",
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("todo")
		},
	}

	clearAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Clear the jumper config and cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("todo")
		},
	}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the jumper config or cache",
	}

	cmd.AddCommand(
		clearAllCmd,
		clearCacheCmd,
		clearConfigCmd,
	)

	return cmd
}
