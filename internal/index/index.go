package index

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/kennycha/ni-idea/internal/store"
)

// Document represents a searchable note document
type Document struct {
	Path    string   `json:"path"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Tags    []string `json:"tags"`
	Content string   `json:"content"`
	Private bool     `json:"private"`
}

// Index wraps bleve index with ni-idea specific operations
type Index struct {
	index bleve.Index
	path  string
}

// SearchOptions configures search behavior
type SearchOptions struct {
	Query          string
	Tags           []string
	Type           store.NoteType
	All            bool
	Limit          int
	IncludePrivate bool
	Fuzzy          bool
	Fuzziness      int // 0 = auto, 1-2 = edit distance
}

// Result represents a search result
type Result struct {
	Note    *store.Note
	Score   float64
	Matches []string
}

// DefaultIndexPath returns the default index path
func DefaultIndexPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "ni-idea", "index"), nil
}

// Open opens an existing index or creates a new one
func Open(indexPath string) (*Index, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return nil, err
	}

	// Try to open existing index
	idx, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// Create new index
		idx, err = bleve.New(indexPath, buildMapping())
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &Index{index: idx, path: indexPath}, nil
}

// Close closes the index
func (i *Index) Close() error {
	return i.index.Close()
}

// IndexNote indexes or updates a note
func (i *Index) IndexNote(note *store.Note) error {
	doc := Document{
		Path:    note.Path,
		Title:   note.Meta.Title,
		Type:    string(note.Meta.Type),
		Tags:    note.Meta.Tags,
		Content: note.Content,
		Private: note.Meta.Private,
	}
	return i.index.Index(note.Path, doc)
}

// DeleteNote removes a note from the index
func (i *Index) DeleteNote(path string) error {
	return i.index.Delete(path)
}

// Search searches the index
func (i *Index) Search(notesDir string, opts SearchOptions) ([]*Result, error) {
	// Build query
	var q query.Query
	if opts.Query == "" {
		q = bleve.NewMatchAllQuery()
	} else if opts.Fuzzy {
		fuzzy := bleve.NewFuzzyQuery(opts.Query)
		if opts.Fuzziness > 0 {
			fuzzy.SetFuzziness(opts.Fuzziness)
		} else {
			fuzzy.SetFuzziness(1)
		}
		q = fuzzy
	} else {
		q = bleve.NewMatchQuery(opts.Query)
	}

	// Create search request
	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = limit
	searchRequest.Fields = []string{"path", "title", "type", "tags", "private"}
	searchRequest.Highlight = bleve.NewHighlight()

	// Execute search
	searchResult, err := i.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Convert to results
	var results []*Result
	for _, hit := range searchResult.Hits {
		path := hit.ID

		// Load full note
		fullPath := filepath.Join(notesDir, path)
		note, err := store.ReadNote(fullPath)
		if err != nil {
			continue
		}
		note.Path = path

		// Apply post-filters
		if !opts.IncludePrivate && note.Meta.Private {
			continue
		}

		// Type filter
		if opts.Type != "" && note.Meta.Type != opts.Type {
			continue
		}

		if !opts.All && opts.Type == "" {
			// Default: only problem and decision
			isDefault := false
			for _, t := range store.DefaultSearchTypes {
				if note.Meta.Type == t {
					isDefault = true
					break
				}
			}
			if !isDefault {
				continue
			}
		}

		// Tag filter
		if len(opts.Tags) > 0 && !hasAnyTag(note.Meta.Tags, opts.Tags) {
			continue
		}

		// Extract match snippets from highlights
		var matches []string
		for _, fragments := range hit.Fragments {
			for _, frag := range fragments {
				// Strip highlight markers
				clean := strings.ReplaceAll(frag, "<mark>", "")
				clean = strings.ReplaceAll(clean, "</mark>", "")
				if len(clean) > 100 {
					clean = clean[:100] + "..."
				}
				matches = append(matches, clean)
			}
		}

		results = append(results, &Result{
			Note:    note,
			Score:   hit.Score,
			Matches: matches,
		})
	}

	return results, nil
}

// Rebuild rebuilds the entire index from notes
func (i *Index) Rebuild(notesDir string) (int, error) {
	// List all notes
	notes, err := store.ListNotes(notesDir, store.ListOptions{IncludePrivate: true})
	if err != nil {
		return 0, err
	}

	// Index each note
	count := 0
	for _, note := range notes {
		if err := i.IndexNote(note); err != nil {
			continue
		}
		count++
	}

	return count, nil
}

// DocCount returns the number of indexed documents
func (i *Index) DocCount() (uint64, error) {
	return i.index.DocCount()
}

// buildMapping creates the index mapping
func buildMapping() mapping.IndexMapping {
	// Text field mapping with English analyzer
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = en.AnalyzerName

	// Keyword field mapping (exact match)
	keywordFieldMapping := bleve.NewKeywordFieldMapping()

	// Boolean field mapping
	boolFieldMapping := bleve.NewBooleanFieldMapping()

	// Document mapping
	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("path", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("tags", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)
	docMapping.AddFieldMappingsAt("private", boolFieldMapping)

	// Index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultAnalyzer = en.AnalyzerName

	return indexMapping
}

func hasAnyTag(noteTags, filterTags []string) bool {
	for _, ft := range filterTags {
		for _, nt := range noteTags {
			if nt == ft {
				return true
			}
		}
	}
	return false
}
