package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init commands for various shells",
	}

	cmd.AddCommand(fishInitCmd())

	return cmd
}

func fishInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fish",
		Short: "Init command for fish shell",
		Args:  cobra.NoArgs,
		RunE:  runFishInitCmd,
	}

	cmd.Flags().String("function", "j", "name of the shell function to define")

	return cmd
}

func runFishInitCmd(cmd *cobra.Command, args []string) error {
	function, err := cmd.Flags().GetString("function")
	if err != nil {
		return err
	}

	sf := fmt.Sprintf(`function %s
    set -l f (mktemp)
    jumper to $argv[1] --out="$f"
    set -l where (cat "$f")
    rm -f "$f"
    cd (realpath "$where"); or return
end;`, function)

	_, err = fmt.Println(sf)

	return err
}
