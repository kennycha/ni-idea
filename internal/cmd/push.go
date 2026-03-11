package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/remote"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var (
	pushRemote         string
	pushAll            bool
	pushIncludePrivate bool
	pushForce          bool
)

var pushCmd = &cobra.Command{
	Use:   "push [path]",
	Short: "Push notes to remote",
	Long: `Push notes to a remote server.

By default, refuses to overwrite if remote version is newer.
Use --force to overwrite regardless of timestamps.

Examples:
  ni push problems/my-note.md           # Push specific note
  ni push --all                         # Push all notes (excludes private)
  ni push --all --include-private       # Push all including private
  ni push --force                       # Force push (ignore conflicts)`,
	RunE: runPush,
}

func init() {
	pushCmd.Flags().StringVar(&pushRemote, "remote", "", "Remote name (uses first remote if not specified)")
	pushCmd.Flags().BoolVar(&pushAll, "all", false, "Push all notes")
	pushCmd.Flags().BoolVar(&pushIncludePrivate, "include-private", false, "Include private notes")
	pushCmd.Flags().BoolVar(&pushForce, "force", false, "Force push (ignore conflicts)")
}

func runPush(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Remotes) == 0 {
		return fmt.Errorf("no remotes configured. Use 'ni remote add' first")
	}

	// Get remote
	var r *config.Remote
	if pushRemote != "" {
		r = cfg.GetRemote(pushRemote)
		if r == nil {
			return fmt.Errorf("remote '%s' not found", pushRemote)
		}
	} else {
		r = cfg.Remotes[0]
	}

	client := remote.NewClient(r.URL, r.Token)

	if pushAll {
		return pushAllNotes(cfg, client)
	}

	if len(args) == 0 {
		return fmt.Errorf("specify a note path or use --all")
	}

	return pushSingleNote(cfg, client, args[0])
}

func pushAllNotes(cfg *config.Config, client *remote.Client) error {
	notes, err := store.ListNotes(cfg.NotesDir, store.ListOptions{
		IncludePrivate: pushIncludePrivate,
	})
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes to push.")
		return nil
	}

	// Get remote notes map for conflict detection
	var remoteNotesMap map[string]*remote.NoteMeta
	if !pushForce {
		remoteNotesMap, _ = client.ListNotesMap()
	}

	fmt.Printf("Pushing %d notes...\n", len(notes))

	pushed := 0
	conflicts := 0
	for _, note := range notes {
		// Skip private notes unless explicitly included
		if note.Meta.Private && !pushIncludePrivate {
			continue
		}

		// Check for conflict (remote is newer)
		if !pushForce && remoteNotesMap != nil {
			if remoteMeta, exists := remoteNotesMap[note.Path]; exists {
				if remoteMeta.Updated > note.Meta.Updated {
					fmt.Printf("  conflict: %s (local: %s, remote: %s)\n",
						note.Path, note.Meta.Updated, remoteMeta.Updated)
					conflicts++
					continue
				}
			}
		}

		// Read full content
		fullPath := filepath.Join(cfg.NotesDir, note.Path)
		content, err := store.ReadNoteFile(fullPath)
		if err != nil {
			fmt.Printf("  skip: %s (failed to read)\n", note.Path)
			continue
		}

		remoteNote := &remote.Note{
			NoteMeta: remote.NoteMeta{
				Path:    note.Path,
				Title:   note.Meta.Title,
				Type:    string(note.Meta.Type),
				Tags:    note.Meta.Tags,
				Private: note.Meta.Private,
				Created: note.Meta.Created,
				Updated: note.Meta.Updated,
			},
			Content: content,
		}

		if err := client.PushNote(remoteNote); err != nil {
			fmt.Printf("  fail: %s (%v)\n", note.Path, err)
			continue
		}

		fmt.Printf("  push: %s\n", note.Path)
		pushed++
	}

	fmt.Printf("Pushed %d/%d notes", pushed, len(notes))
	if conflicts > 0 {
		fmt.Printf(" (%d conflicts)", conflicts)
	}
	fmt.Println()

	if conflicts > 0 {
		fmt.Println("Use --force to overwrite remote versions.")
	}
	return nil
}

func pushSingleNote(cfg *config.Config, client *remote.Client, path string) error {
	// Resolve path
	fullPath, err := store.ResolvePath(cfg.NotesDir, path)
	if err != nil {
		return fmt.Errorf("note not found: %s", path)
	}

	// Read note
	note, err := store.ReadNote(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read note: %w", err)
	}

	// Get relative path
	relPath, _ := filepath.Rel(cfg.NotesDir, fullPath)

	// Check private
	if note.Meta.Private && !pushIncludePrivate {
		return fmt.Errorf("note is private. Use --include-private to push")
	}

	// Check for conflict (remote is newer)
	if !pushForce {
		remoteMeta, err := client.GetNoteMeta(relPath)
		if err == nil && remoteMeta.Updated > note.Meta.Updated {
			return fmt.Errorf("conflict: remote version is newer (local: %s, remote: %s). Use --force to overwrite",
				note.Meta.Updated, remoteMeta.Updated)
		}
	}

	// Read full content
	content, err := store.ReadNoteFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read note content: %w", err)
	}

	remoteNote := &remote.Note{
		NoteMeta: remote.NoteMeta{
			Path:    relPath,
			Title:   note.Meta.Title,
			Type:    string(note.Meta.Type),
			Tags:    note.Meta.Tags,
			Private: note.Meta.Private,
			Created: note.Meta.Created,
			Updated: note.Meta.Updated,
		},
		Content: content,
	}

	if err := client.PushNote(remoteNote); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Printf("Pushed: %s\n", relPath)
	return nil
}
