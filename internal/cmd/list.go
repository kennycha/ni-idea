package cmd

import (
	"fmt"
	"strings"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var (
	listTag            string
	listType           string
	listIncludePrivate bool
)

var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List notes",
	Long: `List notes with optional filtering.

Examples:
  ni list                        # All notes
  ni list problems/              # Notes in specific directory
  ni list --tag infra            # Filter by tag
  ni list --type problem         # Filter by type
  ni list --include-private      # Include private notes`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVar(&listTag, "tag", "", "Filter by tag")
	listCmd.Flags().StringVar(&listType, "type", "", "Filter by type (problem, decision, knowledge, practice)")
	listCmd.Flags().BoolVar(&listIncludePrivate, "include-private", false, "Include private notes")
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	opts := store.ListOptions{
		IncludePrivate: listIncludePrivate,
	}

	if listTag != "" {
		opts.Tags = []string{listTag}
	}

	if listType != "" {
		opts.Type = store.NoteType(listType)
	}

	if len(args) > 0 {
		opts.SubDir = strings.Trim(args[0], "/")
	}

	notes, err := store.ListNotes(cfg.NotesDir, opts)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	for _, note := range notes {
		fmt.Println(note.Path)
		fmt.Printf("  Title: %s\n", note.Meta.Title)
		if len(note.Meta.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(note.Meta.Tags, ", "))
		}
		fmt.Printf("  Updated: %s\n", note.Meta.Updated)
		fmt.Println()
	}

	return nil
}
