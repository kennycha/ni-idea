package store

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kennycha/ni-idea/internal/template"
)

type CreateOptions struct {
	Type         NoteType
	Title        string
	Tags         []string
	Private      bool
	TemplatesDir string
	Body         string
}

func CreateNote(notesDir string, opts CreateOptions) (string, error) {
	// Validate type
	if !opts.Type.IsValid() {
		return "", fmt.Errorf("invalid note type: %s", opts.Type)
	}

	// Get type directory
	typeDir := filepath.Join(notesDir, opts.Type.Directory())
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate filename
	filename := generateFilename(opts.Title)

	// Full path
	notePath := filepath.Join(typeDir, filename)

	// Check if file exists
	if _, err := os.Stat(notePath); err == nil {
		return "", fmt.Errorf("note already exists: %s", notePath)
	}

	// Load template
	content, err := loadTemplate(opts.Type, opts.TemplatesDir)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// Replace placeholders
	today := time.Now().Format("2006-01-02")
	content = strings.ReplaceAll(content, "{{DATE}}", today)

	// Update frontmatter
	content = updateFrontmatter(content, opts)

	// Replace body if provided
	if opts.Body != "" {
		content = replaceBody(content, opts.Body)
	}

	// Write file
	if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write note: %w", err)
	}

	// Return relative path
	relPath, err := filepath.Rel(notesDir, notePath)
	if err != nil {
		return notePath, nil
	}

	return relPath, nil
}

func generateFilename(title string) string {
	if title == "" {
		return time.Now().Format("20060102-150405") + ".md"
	}

	// Slugify: lowercase, replace spaces with hyphens, remove special chars
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters except hyphens and Korean/alphanumeric
	re := regexp.MustCompile(`[^\p{L}\p{N}-]`)
	slug = re.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim hyphens
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return time.Now().Format("20060102-150405") + ".md"
	}

	return slug + ".md"
}

func loadTemplate(noteType NoteType, templatesDir string) (string, error) {
	// Try user template first
	if templatesDir != "" {
		content, err := template.GetTemplate(noteType.String(), templatesDir)
		if err == nil {
			return content, nil
		}
	}

	// Fallback to builtin
	return template.GetBuiltinTemplate(noteType.String()), nil
}

func replaceBody(content string, body string) string {
	lines := strings.Split(content, "\n")
	var result []string
	frontmatterCount := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			frontmatterCount++
		}
		result = append(result, line)
		if frontmatterCount == 2 {
			break
		}
	}

	// Add body after frontmatter
	result = append(result, "")
	result = append(result, body)

	return strings.Join(result, "\n")
}

func updateFrontmatter(content string, opts CreateOptions) string {
	lines := strings.Split(content, "\n")
	var result []string
	inFrontmatter := false
	frontmatterCount := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			frontmatterCount++
			if frontmatterCount == 1 {
				inFrontmatter = true
			} else if frontmatterCount == 2 {
				inFrontmatter = false
			}
			result = append(result, line)
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "title:") && opts.Title != "" {
				result = append(result, fmt.Sprintf("title: \"%s\"", opts.Title))
			} else if strings.HasPrefix(line, "tags:") && len(opts.Tags) > 0 {
				result = append(result, fmt.Sprintf("tags: [%s]", strings.Join(opts.Tags, ", ")))
			} else if strings.HasPrefix(line, "private:") && opts.Private {
				result = append(result, "private: true")
			} else {
				result = append(result, line)
			}
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
