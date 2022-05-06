package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/m-porter/jumper/internal/tui"
	"github.com/spf13/cobra"
)

func ToCmd() *cobra.Command {
	var out string

	cmd := &cobra.Command{
		Use:   "to [query]",
		Args:  cobra.ArbitraryArgs,
		Short: "Display projects in an intractable list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := tui.Run(tui.Options{
				StartingQuery: strings.Join(args, " "),
			})

			if err != nil {
				return err
			}

			switch out {
			case "stderr":
				fmt.Fprintln(os.Stderr, path)
			case "", "stdout":
				fmt.Fprintln(os.Stdout, path)
			default:
				fmt.Fprintln(os.Stdout, path)
				err = writeToFile(out, path)
			}

			return err
		},
	}

	cmd.Flags().StringVar(&out, "out", "stdout", "Where to write the output to")

	return cmd
}

func writeToFile(where string, what string) error {
	f, err := os.OpenFile(where, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(what))
	return err
}
