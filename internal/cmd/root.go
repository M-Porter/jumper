package cmd

import (
	"github.com/spf13/cobra"

	"github.com/m-porter/jumper/internal/logger"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jumper",
		Short: "Seamlessly jump between projects on your machine",
	}

	cmd.PersistentFlags().BoolVar(&logger.Debug, "debug", false, "Run jumper in debug mode")

	cmd.AddCommand(
		ToCmd(),
		AnalyzeCommand(),
		ClearCmd(),
		SetupCmd(),
	)

	return cmd
}
