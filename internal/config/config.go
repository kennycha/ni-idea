package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Remote struct {
	Name  string `yaml:"name" mapstructure:"name"`
	URL   string `yaml:"url" mapstructure:"url"`
	Token string `yaml:"token" mapstructure:"token"`
}

type Config struct {
	NotesDir           string    `yaml:"notes_dir" mapstructure:"notes_dir"`
	TemplatesDir       string    `yaml:"templates_dir" mapstructure:"templates_dir"`
	Editor             string    `yaml:"editor" mapstructure:"editor"`
	DefaultSearchLimit int       `yaml:"default_search_limit" mapstructure:"default_search_limit"`
	Remotes            []*Remote `yaml:"remotes" mapstructure:"remotes"`
}

const (
	ConfigDir  = ".ni-idea"
	ConfigFile = "config.yaml"
)

func ConfigDirPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir)
}

func ConfigFilePath() string {
	return filepath.Join(ConfigDirPath(), ConfigFile)
}

func DefaultNotesDir() string {
	return filepath.Join(ConfigDirPath(), "notes")
}

func DefaultTemplatesDir() string {
	return filepath.Join(ConfigDirPath(), "templates")
}

func DefaultConfig() *Config {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return &Config{
		NotesDir:           DefaultNotesDir(),
		TemplatesDir:       DefaultTemplatesDir(),
		Editor:             editor,
		DefaultSearchLimit: 10,
	}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigFile(ConfigFilePath())
	viper.SetConfigType("yaml")

	// Environment variable bindings
	viper.SetEnvPrefix("NI")
	viper.BindEnv("notes_dir", "NI_NOTES_DIR")
	viper.BindEnv("templates_dir", "NI_TEMPLATES_DIR")
	viper.BindEnv("editor", "NI_EDITOR")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file exists but has errors
			return nil, err
		}
		// Config file not found, use defaults
		return cfg, nil
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Expand ~ in paths
	cfg.NotesDir = expandPath(cfg.NotesDir)
	cfg.TemplatesDir = expandPath(cfg.TemplatesDir)

	return cfg, nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}

func simplifyPath(path string) string {
	home, _ := os.UserHomeDir()
	if filepath.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

func (c *Config) ToYAML() string {
	yaml := `notes_dir: ` + simplifyPath(c.NotesDir) + `
templates_dir: ` + simplifyPath(c.TemplatesDir) + `
default_search_limit: ` + fmt.Sprintf("%d", c.DefaultSearchLimit) + `
`
	if len(c.Remotes) > 0 {
		yaml += "remotes:\n"
		for _, r := range c.Remotes {
			yaml += fmt.Sprintf("  - name: %s\n    url: %s\n    token: %s\n", r.Name, r.URL, r.Token)
		}
	}
	return yaml
}

// GetRemote returns a remote by name
func (c *Config) GetRemote(name string) *Remote {
	for _, r := range c.Remotes {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// AddRemote adds a new remote
func (c *Config) AddRemote(remote *Remote) error {
	if c.GetRemote(remote.Name) != nil {
		return fmt.Errorf("remote '%s' already exists", remote.Name)
	}
	c.Remotes = append(c.Remotes, remote)
	return nil
}

// RemoveRemote removes a remote by name
func (c *Config) RemoveRemote(name string) error {
	for i, r := range c.Remotes {
		if r.Name == name {
			c.Remotes = append(c.Remotes[:i], c.Remotes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("remote '%s' not found", name)
}

// Save writes the config to the config file
func (c *Config) Save() error {
	return os.WriteFile(ConfigFilePath(), []byte(c.ToYAML()), 0644)
}
