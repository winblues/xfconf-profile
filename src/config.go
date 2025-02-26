package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed assets/config.yml
var defaultConfig []byte

type MergeBehavior string
type ExcludePatterns map[string]*regexp.Regexp

const (
	MergeSoft  MergeBehavior = "soft"
	MergeHard  MergeBehavior = "hard"
	MergeForce MergeBehavior = "force"
)

func (ep *ExcludePatterns) IsExcluded(channel string, property string) bool {
	representation := fmt.Sprintf("%s%s", channel, property)
	for _, re := range *ep {
		if re.MatchString(representation) {
			return true
		}
	}
	return false
}

type Config struct {
	Version int `yaml:"version"`
	Sync    struct {
		Auto bool `yaml:"auto"`
	} `yaml:"sync"`
	Merge   MergeBehavior   `yaml:"merge"`
	Exclude ExcludePatterns `yaml:"exclude"`
}

func ParseMergeBehavior(value string) (MergeBehavior, error) {
	switch strings.ToLower(value) {
	case "soft":
		return MergeSoft, nil
	case "hard":
		return MergeHard, nil
	case "force":
		return MergeForce, nil
	default:
		return "", errors.New("invalid merge behavior: must be 'soft', 'hard', or 'force'")
	}
}

func (m *MergeBehavior) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	parsed, err := ParseMergeBehavior(raw)
	if err != nil {
		return err
	}

	*m = parsed
	return nil
}

func (m *ExcludePatterns) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var patternStrings []string
	if err := unmarshal(&patternStrings); err != nil {
		return err
	}

	patterns := make(ExcludePatterns)

	for _, pattern := range patternStrings {
		// Compile the regular expression for each pattern
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid exclude regular expression: %s", pattern)
		}

		patterns[pattern] = re
	}

	*m = patterns
	return nil
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
