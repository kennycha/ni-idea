package cmd

import (
	"fmt"
	"os"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get full content of a note",
	Long: `Get the full content of a note by its path.
The path is relative to the notes directory.

Examples:
  ni get problems/frontend/nextjs-caching
  ni get knowledge/infra/k8s-basics.md`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	notePath := args[0]

	// Resolve the actual file path
	fullPath, err := store.ResolvePath(cfg.NotesDir, notePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Note not found: %s\n", notePath)
		return nil
	}

	// Read and print the file content
	content, err := store.ReadNoteFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read note: %w", err)
	}

	fmt.Print(content)
	return nil
}
