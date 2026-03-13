package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/index"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var (
	addType    string
	addTitle   string
	addTags    string
	addPrivate bool
	addNoEdit  bool
	addBody    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new note",
	Long: `Add a new note with the specified type.

Examples:
  ni add --type problem --title "캐싱 이슈"
  ni add --type decision --tag infra,k8s
  ni add --type problem --private
  ni add --type problem --no-edit`,
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addType, "type", "", "Note type (problem, decision, knowledge, practice) [required]")
	addCmd.Flags().StringVar(&addTitle, "title", "", "Note title")
	addCmd.Flags().StringVar(&addTags, "tag", "", "Tags (comma separated)")
	addCmd.Flags().BoolVar(&addPrivate, "private", false, "Mark as private")
	addCmd.Flags().BoolVar(&addNoEdit, "no-edit", false, "Don't open editor")
	addCmd.Flags().StringVar(&addBody, "body", "", "Note body content (skips editor)")
	addCmd.MarkFlagRequired("type")
}

func runAdd(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate type
	noteType := store.NoteType(addType)
	if !noteType.IsValid() {
		return fmt.Errorf("invalid type: %s (must be problem, decision, knowledge, or practice)", addType)
	}

	// Parse tags
	var tags []string
	if addTags != "" {
		for _, tag := range strings.Split(addTags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	// Create note
	opts := store.CreateOptions{
		Type:         noteType,
		Title:        addTitle,
		Tags:         tags,
		Private:      addPrivate,
		TemplatesDir: cfg.TemplatesDir,
		Body:         addBody,
	}

	notePath, err := store.CreateNote(cfg.NotesDir, opts)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("Created: %s\n", notePath)

	fullPath := filepath.Join(cfg.NotesDir, notePath)

	// Open in editor (skip if --body provided or --no-edit)
	if !addNoEdit && addBody == "" {
		editor := cfg.Editor
		if editor == "" {
			editor = os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
		}

		editorCmd := exec.Command(editor, fullPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}
	}

	// Index the note
	if err := indexNote(fullPath, notePath); err != nil {
		fmt.Printf("Warning: failed to index note: %v\n", err)
	}

	return nil
}

func indexNote(fullPath, relativePath string) error {
	indexPath, err := index.DefaultIndexPath()
	if err != nil {
		return err
	}

	idx, err := index.Open(indexPath)
	if err != nil {
		return err
	}
	defer idx.Close()

	note, err := store.ReadNote(fullPath)
	if err != nil {
		return err
	}
	note.Path = relativePath

	return idx.IndexNote(note)
}
