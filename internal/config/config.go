package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	NotesDir           string `yaml:"notes_dir" mapstructure:"notes_dir"`
	TemplatesDir       string `yaml:"templates_dir" mapstructure:"templates_dir"`
	Editor             string `yaml:"editor" mapstructure:"editor"`
	DefaultSearchLimit int    `yaml:"default_search_limit" mapstructure:"default_search_limit"`
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
	return `notes_dir: ` + simplifyPath(c.NotesDir) + `
templates_dir: ` + simplifyPath(c.TemplatesDir) + `
default_search_limit: ` + fmt.Sprintf("%d", c.DefaultSearchLimit) + `
`
}
