package cmd

import (
	"fmt"
	"strings"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/index"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var (
	searchTag            string
	searchType           string
	searchAll            bool
	searchLimit          int
	searchIncludePrivate bool
	searchFuzzy          bool
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes by keyword",
	Long: `Search notes by keyword with optional filtering.
By default, only searches problem and decision types.

Examples:
  ni search "Next.js 캐싱"
  ni search "캐싱" --tag nextjs
  ni search "k8s" --type problem
  ni search "캐싱" --all`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().StringVar(&searchTag, "tag", "", "Filter by tag")
	searchCmd.Flags().StringVar(&searchType, "type", "", "Filter by type (problem, decision, knowledge, practice)")
	searchCmd.Flags().BoolVar(&searchAll, "all", false, "Search all note types")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 0, "Limit number of results (0 = use config default)")
	searchCmd.Flags().BoolVar(&searchIncludePrivate, "include-private", false, "Include private notes")
	searchCmd.Flags().BoolVar(&searchFuzzy, "fuzzy", false, "Enable fuzzy search (allows typos)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	query := args[0]

	// Determine limit
	limit := searchLimit
	if limit == 0 {
		limit = cfg.DefaultSearchLimit
	}

	// Open index
	indexPath, err := index.DefaultIndexPath()
	if err != nil {
		return fmt.Errorf("failed to get index path: %w", err)
	}

	idx, err := index.Open(indexPath)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer idx.Close()

	// Check if index is empty and rebuild if needed
	docCount, _ := idx.DocCount()
	if docCount == 0 {
		fmt.Println("Index is empty, building...")
		count, err := idx.Rebuild(cfg.NotesDir)
		if err != nil {
			return fmt.Errorf("failed to build index: %w", err)
		}
		fmt.Printf("Indexed %d notes.\n\n", count)
	}

	opts := index.SearchOptions{
		Query:          query,
		All:            searchAll,
		Limit:          limit,
		IncludePrivate: searchIncludePrivate,
		Fuzzy:          searchFuzzy,
	}

	if searchTag != "" {
		opts.Tags = []string{searchTag}
	}

	if searchType != "" {
		opts.Type = store.NoteType(searchType)
	}

	results, err := idx.Search(cfg.NotesDir, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	for _, result := range results {
		fmt.Println(result.Note.Path)
		fmt.Printf("  Title: %s\n", result.Note.Meta.Title)
		if len(result.Note.Meta.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(result.Note.Meta.Tags, ", "))
		}
		if len(result.Matches) > 0 {
			// Show first match as snippet
			fmt.Printf("  Match: \"%s\"\n", result.Matches[0])
		}
		fmt.Println()
	}

	return nil
}
