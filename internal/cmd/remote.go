package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/kennycha/ni-idea/internal/config"
	"github.com/kennycha/ni-idea/internal/remote"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote servers",
	Long: `Manage remote servers for syncing notes.

Examples:
  ni remote add personal https://my-server.com
  ni remote list
  ni remote remove personal`,
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new remote",
	Long:  `Add a new remote server for syncing notes.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runRemoteAdd,
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured remotes",
	RunE:  runRemoteList,
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a remote",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemoteRemove,
}

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)
}

func runRemoteAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	url := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if remote already exists
	if cfg.GetRemote(name) != nil {
		return fmt.Errorf("remote '%s' already exists", name)
	}

	// Prompt for token
	fmt.Print("Token: ")
	token, err := readPassword()
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	fmt.Println()

	if token == "" {
		return fmt.Errorf("token is required")
	}

	// Test connection
	fmt.Print("Testing connection... ")
	client := remote.NewClient(url, token)
	if err := client.Ping(); err != nil {
		fmt.Println("failed")
		return fmt.Errorf("connection test failed: %w", err)
	}
	fmt.Println("ok")

	// Add remote
	r := &config.Remote{
		Name:  name,
		URL:   url,
		Token: token,
	}

	if err := cfg.AddRemote(r); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Added remote '%s' (%s)\n", name, url)
	return nil
}

func runRemoteList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Remotes) == 0 {
		fmt.Println("No remotes configured.")
		fmt.Println("Use 'ni remote add <name> <url>' to add one.")
		return nil
	}

	for _, r := range cfg.Remotes {
		fmt.Printf("%s\t%s\n", r.Name, r.URL)
	}

	return nil
}

func runRemoteRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.RemoveRemote(name); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Removed remote '%s'\n", name)
	return nil
}

func readPassword() (string, error) {
	// Check if stdin is a terminal
	if term.IsTerminal(int(syscall.Stdin)) {
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		return string(password), nil
	}

	// If not a terminal, read from stdin directly
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}
