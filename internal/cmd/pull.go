package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/index"
	"github.com/kennycha/ni-idea/internal/remote"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/spf13/cobra"
)

var (
	pullRemote string
	pullForce  bool
	pullOurs   bool
	pullTheirs bool
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull notes from remote",
	Long: `Pull notes from a remote server.

By default, only downloads notes that don't exist locally.
When a note exists both locally and remotely with different timestamps,
it's treated as a conflict.

Examples:
  ni pull                    # Pull from default remote (skip conflicts)
  ni pull --force            # Overwrite all existing notes
  ni pull --theirs           # Use remote version for conflicts
  ni pull --ours             # Keep local version for conflicts (skip)`,
	RunE: runPull,
}

func init() {
	pullCmd.Flags().StringVar(&pullRemote, "remote", "", "Remote name (uses first remote if not specified)")
	pullCmd.Flags().BoolVar(&pullForce, "force", false, "Overwrite all existing notes")
	pullCmd.Flags().BoolVar(&pullTheirs, "theirs", false, "Use remote version for conflicts")
	pullCmd.Flags().BoolVar(&pullOurs, "ours", false, "Keep local version for conflicts")
}

func runPull(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Remotes) == 0 {
		return fmt.Errorf("no remotes configured. Use 'ni remote add' first")
	}

	// Get remote
	var r *config.Remote
	if pullRemote != "" {
		r = cfg.GetRemote(pullRemote)
		if r == nil {
			return fmt.Errorf("remote '%s' not found", pullRemote)
		}
	} else {
		r = cfg.Remotes[0]
	}

	client := remote.NewClient(r.URL, r.Token)

	// List remote notes
	fmt.Printf("Fetching notes from '%s'...\n", r.Name)
	remoteMetas, err := client.ListNotes()
	if err != nil {
		return fmt.Errorf("failed to list remote notes: %w", err)
	}

	if len(remoteMetas) == 0 {
		fmt.Println("No notes on remote.")
		return nil
	}

	pulled := 0
	skipped := 0
	conflicts := 0

	for _, meta := range remoteMetas {
		localPath := filepath.Join(cfg.NotesDir, meta.Path)

		// Check if exists locally
		if _, err := os.Stat(localPath); err == nil {
			// File exists locally - check for conflict
			if pullForce || pullTheirs {
				// Force overwrite
			} else if pullOurs {
				skipped++
				continue
			} else {
				// Check timestamps for conflict
				localNote, err := store.ReadNote(localPath)
				if err == nil && localNote.Meta.Updated != meta.Updated {
					// Conflict detected
					fmt.Printf("  conflict: %s (local: %s, remote: %s)\n",
						meta.Path, localNote.Meta.Updated, meta.Updated)
					conflicts++
					continue
				}
				// Same timestamp or can't read local - skip
				skipped++
				continue
			}
		}

		// Fetch full note
		note, err := client.GetNote(meta.Path)
		if err != nil {
			fmt.Printf("  fail: %s (%v)\n", meta.Path, err)
			continue
		}

		// Ensure directory exists
		dir := filepath.Dir(localPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("  fail: %s (cannot create directory)\n", meta.Path)
			continue
		}

		// Write note
		if err := os.WriteFile(localPath, []byte(note.Content), 0644); err != nil {
			fmt.Printf("  fail: %s (cannot write file)\n", meta.Path)
			continue
		}

		fmt.Printf("  pull: %s\n", meta.Path)
		pulled++
	}

	fmt.Printf("Pulled %d notes", pulled)
	if conflicts > 0 {
		fmt.Printf(" (%d conflicts", conflicts)
		if skipped > 0 {
			fmt.Printf(", %d skipped", skipped)
		}
		fmt.Print(")")
	} else if skipped > 0 {
		fmt.Printf(" (skipped %d existing)", skipped)
	}
	fmt.Println()

	if conflicts > 0 {
		fmt.Println("Use --theirs to overwrite with remote version, or --ours to keep local.")
	}

	// Rebuild index if any notes were pulled
	if pulled > 0 {
		fmt.Print("Rebuilding index... ")
		indexPath, err := index.DefaultIndexPath()
		if err == nil {
			idx, err := index.Open(indexPath)
			if err == nil {
				idx.Rebuild(cfg.NotesDir)
				idx.Close()
				fmt.Println("done")
			}
		}
	}

	return nil
}
