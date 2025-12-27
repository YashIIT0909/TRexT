package storage

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	Theme          string `yaml:"theme"`
	DefaultTimeout int    `yaml:"defaultTimeout"` // seconds
	SSLVerify      bool   `yaml:"sslVerify"`
	Proxy          string `yaml:"proxy"`
	History        struct {
		MaxItems int  `yaml:"maxItems"`
		Enabled  bool `yaml:"enabled"`
	} `yaml:"history"`
	Keybindings struct {
		SendRequest string `yaml:"sendRequest"`
		NewRequest  string `yaml:"newRequest"`
		SaveRequest string `yaml:"saveRequest"`
		FocusURL    string `yaml:"focusURL"`
	} `yaml:"keybindings"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{
		Theme:          "default",
		DefaultTimeout: 30,
		SSLVerify:      true,
	}
	cfg.History.MaxItems = 100
	cfg.History.Enabled = true
	cfg.Keybindings.SendRequest = "Ctrl+Enter"
	cfg.Keybindings.NewRequest = "Ctrl+N"
	cfg.Keybindings.SaveRequest = "Ctrl+S"
	cfg.Keybindings.FocusURL = "Ctrl+U"
	return cfg
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "trext", "config.yaml"), nil
}

// LoadConfig loads configuration from file or returns default
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			_ = SaveConfig(cfg) // Try to save default config
			return cfg, nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
