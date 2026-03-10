package store

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ListOptions struct {
	IncludePrivate bool
	Tags           []string
	Type           NoteType
	SubDir         string // 특정 서브디렉토리만 검색
}

// ReadNote reads a note file and parses frontmatter + content
func ReadNote(path string) (*Note, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	note, err := parseNote(string(content))
	if err != nil {
		return nil, err
	}
	note.Path = path

	return note, nil
}

// ReadNoteFile reads raw file content (for ni get)
func ReadNoteFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ResolvePath finds the actual file path for a note
// Handles both with and without .md extension
func ResolvePath(notesDir, notePath string) (string, error) {
	// Remove leading/trailing slashes
	notePath = strings.Trim(notePath, "/")

	// Try exact path first
	fullPath := filepath.Join(notesDir, notePath)
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, nil
	}

	// Try with .md extension
	if !strings.HasSuffix(notePath, ".md") {
		fullPathMd := fullPath + ".md"
		if _, err := os.Stat(fullPathMd); err == nil {
			return fullPathMd, nil
		}
	}

	return "", os.ErrNotExist
}

// ListNotes lists all notes in the notes directory with optional filtering
func ListNotes(notesDir string, opts ListOptions) ([]*Note, error) {
	var notes []*Note

	searchDir := notesDir
	if opts.SubDir != "" {
		searchDir = filepath.Join(notesDir, opts.SubDir)
	}

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip directories we can't read
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		note, err := ReadNote(path)
		if err != nil {
			return nil // Skip files we can't parse
		}

		// Make path relative to notesDir
		relPath, _ := filepath.Rel(notesDir, path)
		note.Path = relPath

		// Apply filters
		if !opts.IncludePrivate && note.Meta.Private {
			return nil
		}

		if opts.Type != "" && note.Meta.Type != opts.Type {
			return nil
		}

		if len(opts.Tags) > 0 && !hasAnyTag(note.Meta.Tags, opts.Tags) {
			return nil
		}

		notes = append(notes, note)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return notes, nil
}

// CollectTags collects all tags with their counts
func CollectTags(notesDir string) (map[string]int, error) {
	tags := make(map[string]int)

	notes, err := ListNotes(notesDir, ListOptions{IncludePrivate: false})
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		for _, tag := range note.Meta.Tags {
			tags[tag]++
		}
	}

	return tags, nil
}

// parseNote parses frontmatter and content from a note string
func parseNote(content string) (*Note, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var frontmatter strings.Builder
	var body strings.Builder
	inFrontmatter := false
	frontmatterDone := false
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if lineNum == 1 && line == "---" {
			inFrontmatter = true
			continue
		}

		if inFrontmatter && line == "---" {
			inFrontmatter = false
			frontmatterDone = true
			continue
		}

		if inFrontmatter {
			frontmatter.WriteString(line)
			frontmatter.WriteString("\n")
		} else if frontmatterDone {
			body.WriteString(line)
			body.WriteString("\n")
		}
	}

	var meta NoteMeta
	if err := yaml.Unmarshal([]byte(frontmatter.String()), &meta); err != nil {
		return nil, err
	}

	return &Note{
		Meta:    meta,
		Content: strings.TrimSpace(body.String()),
	}, nil
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
