package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/m-porter/jumper/internal/logger"
)

func RootCmd() *cobra.Command {
	version := "debug"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		version = info.Main.Version
	}

	cmd := &cobra.Command{
		Use:     "jumper",
		Short:   "Seamlessly jump between projects on your machine",
		Version: version,
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
