package cmd

import (
	"fmt"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/index"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Manage search index",
	Long: `Manage the search index for faster queries.

Examples:
  ni index rebuild    # Rebuild the entire index
  ni index status     # Show index status`,
}

var indexRebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild the search index",
	Long:  `Rebuild the entire search index from all notes.`,
	RunE:  runIndexRebuild,
}

var indexStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show index status",
	Long:  `Show the current status of the search index.`,
	RunE:  runIndexStatus,
}

func init() {
	indexCmd.AddCommand(indexRebuildCmd)
	indexCmd.AddCommand(indexStatusCmd)
}

func runIndexRebuild(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	indexPath, err := index.DefaultIndexPath()
	if err != nil {
		return fmt.Errorf("failed to get index path: %w", err)
	}

	fmt.Println("Rebuilding index...")

	idx, err := index.Open(indexPath)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer idx.Close()

	count, err := idx.Rebuild(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed to rebuild index: %w", err)
	}

	fmt.Printf("Indexed %d notes.\n", count)
	return nil
}

func runIndexStatus(cmd *cobra.Command, args []string) error {
	indexPath, err := index.DefaultIndexPath()
	if err != nil {
		return fmt.Errorf("failed to get index path: %w", err)
	}

	idx, err := index.Open(indexPath)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer idx.Close()

	count, err := idx.DocCount()
	if err != nil {
		return fmt.Errorf("failed to get doc count: %w", err)
	}

	fmt.Printf("Index path: %s\n", indexPath)
	fmt.Printf("Indexed documents: %d\n", count)
	return nil
}
