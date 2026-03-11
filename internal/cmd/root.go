package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ni",
	Short: "Personal knowledge base CLI tool",
	Long: `ni-idea is a personal knowledge base CLI tool.
It stores knowledge locally in markdown format and provides fast search.

Examples:
  ni init                          # Initialize notes directory and config
  ni search "query"                # Search problems/decisions
  ni search "query" --all          # Search all note types
  ni get problems/frontend/caching # Get full note content
  ni list                          # List all notes
  ni add                           # Add a new note`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
}
