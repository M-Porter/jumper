package cmd

import "github.com/spf13/cobra"

var (
	runInDebugMode bool
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jumper",
		Short: "Seamlessly jump between projects on your machine",
	}

	cmd.PersistentFlags().BoolVar(&runInDebugMode, "debug", false, "Run jumper in debug mode.")

	cmd.AddCommand(
		ToCmd(),
		AnalyzeCommand(),
		ClearCmd(),
		InstallCmd(),
	)

	return cmd
}