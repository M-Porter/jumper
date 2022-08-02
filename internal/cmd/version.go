package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func VersionCmd(options RootCmdOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print jumper version",
		Run: func(cmd *cobra.Command, args []string) {
			var prettyDate = options.Date
			parsedDate, err := time.Parse(time.RFC3339, options.Date)
			if err == nil {
				prettyDate = parsedDate.Format(time.UnixDate)
			}

			fmt.Println("Jumper:")
			fmt.Printf("  Version:    %s\n", options.Version)
			fmt.Printf("  Commit:     %s\n", options.Commit)
			fmt.Printf("  Built by:   %s\n", options.BuiltBy)
			fmt.Printf("  Build date: %s\n", prettyDate)
		},
	}
}
