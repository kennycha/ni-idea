package cmd

import (
	"fmt"
	"sort"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags with note counts",
	Long: `List all tags used in notes with their counts.

Examples:
  ni tags`,
	RunE: runTags,
}

type tagCount struct {
	Tag   string
	Count int
}

func runTags(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	tags, err := store.CollectTags(cfg.NotesDir)
	if err != nil {
		return fmt.Errorf("failed to collect tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found.")
		return nil
	}

	// Sort by count (descending), then by name
	var sorted []tagCount
	for tag, count := range tags {
		sorted = append(sorted, tagCount{Tag: tag, Count: count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Count != sorted[j].Count {
			return sorted[i].Count > sorted[j].Count
		}
		return sorted[i].Tag < sorted[j].Tag
	})

	// Find max tag length for alignment
	maxLen := 0
	for _, tc := range sorted {
		if len(tc.Tag) > maxLen {
			maxLen = len(tc.Tag)
		}
	}

	for _, tc := range sorted {
		fmt.Printf("%-*s  (%d)\n", maxLen, tc.Tag, tc.Count)
	}

	return nil
}
