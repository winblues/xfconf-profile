package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed assets/config.yml
var defaultConfig []byte

type Config struct {
	Version int `yaml:"version"`
	Sync    struct {
		Enable bool `yaml:"enable"`
	} `yaml:"sync"`
	Merge   struct{} `yaml:"merge"`
	Exclude []string `yaml:"exclude"`
}

func getConfigPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		xdgConfigHome = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgConfigHome, "xfconf-profile", "config.yml")
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write default config if it does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
