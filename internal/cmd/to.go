package cmd

import (
	"fmt"
	"strings"

	"github.com/m-porter/jumper/internal/tui2"

	"github.com/spf13/cobra"
)

func ToCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to [query]",
		Args:  cobra.ArbitraryArgs,
		Short: "Display projects in an intractable list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			//path, err := tui.Run(runInDebugMode, strings.Join(args, " "))
			path, err := tui2.Run(runInDebugMode, strings.Join(args, " "))

			if err != nil {
				return err
			}

			fmt.Println(path)

			return nil
		},
	}
}
