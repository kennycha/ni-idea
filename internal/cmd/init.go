package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/skill"
	"github.com/kennycha/ni-idea/internal/store"
	"github.com/kennycha/ni-idea/internal/template"
	"github.com/spf13/cobra"
)

var (
	initNotesDir string
	initForce    bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ni-idea configuration and notes directory",
	Long: `Initialize ni-idea by creating:
- Configuration file at ~/.ni-idea/config.yaml
- Notes directory structure at ~/.ni-idea/notes/
- Note templates at ~/.ni-idea/templates/

Examples:
  ni init
  ni init --notes-dir ~/my-notes
  ni init --force`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&initNotesDir, "notes-dir", "", "Notes directory path (default: ~/.ni-idea/notes)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing configuration")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine notes directory
	notesDir := initNotesDir
	if notesDir == "" {
		notesDir = config.DefaultNotesDir()
	} else if notesDir[0] == '~' {
		home, _ := os.UserHomeDir()
		notesDir = filepath.Join(home, notesDir[1:])
	}

	// Check existing config
	configPath := config.ConfigFilePath()
	if _, err := os.Stat(configPath); err == nil && !initForce {
		fmt.Printf("Config already exists: %s\n", configPath)
		fmt.Println("Use --force to overwrite")
		return nil
	}

	// Create config directory
	configDir := config.ConfigDirPath()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config file
	cfg := config.DefaultConfig()
	cfg.NotesDir = notesDir
	if err := os.WriteFile(configPath, []byte(cfg.ToYAML()), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	fmt.Printf("Created config: %s\n", configPath)

	// Create notes directory structure
	if err := createNotesDirectories(notesDir); err != nil {
		return err
	}
	fmt.Printf("Created notes directory: %s\n", notesDir)

	// Create templates
	templatesDir := config.DefaultTemplatesDir()
	if err := template.WriteTemplates(templatesDir); err != nil {
		return fmt.Errorf("failed to create templates: %w", err)
	}
	fmt.Printf("Created templates: %s\n", templatesDir)

	// Ask about skill installation
	if shouldInstallSkill() {
		if err := skill.InstallSkills(); err != nil {
			fmt.Printf("Warning: failed to install skills: %v\n", err)
		} else {
			fmt.Println("Installed Claude Code skills:")
			fmt.Println("  - ~/.claude/skills/ni-idea-search/")
			fmt.Println("  - ~/.claude/skills/ni-idea-add/")
		}
	}

	return nil
}

func shouldInstallSkill() bool {
	fmt.Print("Install Claude Code skills? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func createNotesDirectories(notesDir string) error {
	// Create base notes directory
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %w", err)
	}

	// Create type directories
	for _, noteType := range store.AllTypes {
		dir := filepath.Join(notesDir, noteType.Directory())
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", noteType.Directory(), err)
		}
		fmt.Printf("  - %s/\n", noteType.Directory())
	}

	return nil
}
