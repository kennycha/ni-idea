package search

import (
	"strings"

	"github.com/kennycha/ni-idea/internal/store"
)

type Options struct {
	Query          string
	Tags           []string
	Type           store.NoteType
	All            bool // true면 모든 타입 검색
	Limit          int
	IncludePrivate bool
}

type Result struct {
	Note    *store.Note
	Matches []string // 매칭된 라인들
}

// Search searches notes by query with optional filtering
func Search(notesDir string, opts Options) ([]*Result, error) {
	// Determine which types to search
	var searchTypes []store.NoteType
	if opts.Type != "" {
		searchTypes = []store.NoteType{opts.Type}
	} else if opts.All {
		searchTypes = store.AllTypes
	} else {
		searchTypes = store.DefaultSearchTypes
	}

	var results []*Result

	for _, noteType := range searchTypes {
		listOpts := store.ListOptions{
			IncludePrivate: opts.IncludePrivate,
			Tags:           opts.Tags,
			Type:           noteType,
		}

		notes, err := store.ListNotes(notesDir, listOpts)
		if err != nil {
			continue
		}

		for _, note := range notes {
			matches := searchInNote(note, opts.Query)
			if len(matches) > 0 || matchesFilename(note.Path, opts.Query) {
				results = append(results, &Result{
					Note:    note,
					Matches: matches,
				})
			}

			// Apply limit
			if opts.Limit > 0 && len(results) >= opts.Limit {
				return results, nil
			}
		}
	}

	return results, nil
}

// searchInNote searches for query in note content and title
func searchInNote(note *store.Note, query string) []string {
	var matches []string
	queryLower := strings.ToLower(query)

	// Search in title
	if strings.Contains(strings.ToLower(note.Meta.Title), queryLower) {
		matches = append(matches, "Title: "+note.Meta.Title)
	}

	// Search in content
	lines := strings.Split(note.Content, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), queryLower) {
			// Trim and limit line length
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 100 {
				trimmed = trimmed[:100] + "..."
			}
			if trimmed != "" {
				matches = append(matches, trimmed)
			}
		}
	}

	return matches
}

// matchesFilename checks if query matches the filename
func matchesFilename(path, query string) bool {
	return strings.Contains(strings.ToLower(path), strings.ToLower(query))
}
