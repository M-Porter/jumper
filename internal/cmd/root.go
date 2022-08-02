package cmd

import (
	"github.com/spf13/cobra"

	"github.com/m-porter/jumper/internal/logger"
)

type RootCmdOptions struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func RootCmd(options RootCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jumper",
		Short:   "Seamlessly jump between projects on your machine",
		Version: options.Version,
	}

	cmd.PersistentFlags().BoolVar(&logger.Debug, "debug", false, "Run jumper in debug mode")

	cmd.AddCommand(
		ToCmd(),
		AnalyzeCommand(),
		ClearCmd(),
		VersionCmd(options),
		EditCmd(),
	)

	return cmd
}
