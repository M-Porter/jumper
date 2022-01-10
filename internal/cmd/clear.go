package cmd

import (
	"fmt"
	"os"

	"github.com/m-porter/jumper/internal/config"

	"github.com/spf13/cobra"
)

func ClearCmd() *cobra.Command {
	clearCacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Delete the jumper cache",
		Run: func(cmd *cobra.Command, args []string) {
			clearCache()
		},
	}

	clearConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Delete the jumper config",
		Run: func(cmd *cobra.Command, args []string) {
			clearConfig()
		},
	}

	clearAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Delete the jumper config and cache",
		Run: func(cmd *cobra.Command, args []string) {
			clearCache()
			clearConfig()
		},
	}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Delete the jumper config or cache",
	}

	cmd.AddCommand(
		clearAllCmd,
		clearCacheCmd,
		clearConfigCmd,
	)

	return cmd
}

func clearCache() {
	deleteFile(config.C.CacheFileFullPath, "Deleting cache", "Could not delete cache")
}

func clearConfig() {
	deleteFile(config.Filepath(), "Deleting config", "Could not delete config")
}

func deleteFile(path, actionMsg, deleteFileErrorMsg string) {
	fmt.Printf("%s...", actionMsg)

	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			fmt.Printf("Error! %s\n", deleteFileErrorMsg)
			return
		}
	}

	fmt.Println("Done")
}
