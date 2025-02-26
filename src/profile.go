package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Profile map[string]map[string]any

// TODO: return a new profile that only includes properties that were actually changed based on the merge and exclude settings.
func applyProfile(profilePath string, mergeBehavior MergeBehavior, exclude ExcludePatterns, dryRun bool) error {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var profile Profile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	for channel, properties := range profile {
		// Keys starting with X- are not channels
		if strings.HasPrefix(channel, "X-") {
			continue
		}

		for property, value := range properties {
			// Check if this property should be skipped based on merge preferences
			if mergeBehavior == MergeSoft {
				currentValue := "currentValue"
				defaultValue := "defaultValue"

				if currentValue != defaultValue {
					fmt.Printf("%s Skipping property %s with non-default value %s (default=%s)%s\n", yellow("•"), channel, property, currentValue, defaultValue)
					continue
				}
			}

			if exclude.IsExcluded(channel, property) && mergeBehavior != MergeForce {
				fmt.Printf("%s Skipping excluded property %s%s\n", yellow("•"), channel, property)
				continue
			}

			dryRunNotice := ""
			if dryRun {
				dryRunNotice = " (skipping due to dry run)"
			}

			fmt.Printf("%s Setting %s%s ➔ %s%s\n", blue("•"), channel, property, value, dryRunNotice)

			if dryRun {
				continue
			}

			// We can definitely set this property now

			cmd := exec.Command("xfconf-query", "-c", channel, "--property", property, "--type", "string", "--create", "--set", fmt.Sprintf("%v", value))

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
			}
		}
	}

	return nil
}

func revertProfile(profilePath string, exclude ExcludePatterns, dryRun bool) error {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var profile Profile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	for channel, properties := range profile {
		// Keys starting with X- are not channels
		if strings.HasPrefix(channel, "X-") {
			continue
		}

		for property := range properties {
			if exclude.IsExcluded(channel, property) {
				fmt.Printf("%s Skipping excluded property %s%s\n", yellow("•"), channel, property)
				continue
			}

			dryRunNotice := ""
			if dryRun {
				dryRunNotice = " (skipping due to dry run)"
			}

			fmt.Printf("%s Resetting %s%s%s\n", blue("•"), channel, property, dryRunNotice)

			if dryRun {
				continue
			}

			// We can definitely reset this property now
			cmd := exec.Command("xfconf-query", "-c", channel, "--reset", "--property", property)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
			}
		}
	}

	return nil
}

// Create $XDG_STATE_HOME/xfconf-profile/sync if needed
func ensureStateDir() (string, error) {
	xdgStateHome := os.Getenv("XDG_STATE_HOME")
	var stateDirPath string

	if xdgStateHome != "" {
		stateDirPath = filepath.Join(xdgStateHome, "xfconf-profile", "sync")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		stateDirPath = filepath.Join(homeDir, ".local", "state", "xfconf-profile", "sync")
	}

	if err := os.MkdirAll(stateDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create state directory: %v", err)
	}

	return stateDirPath, nil
}

func copyDistConfig(distConfig string, currentDir string) error {
	defaultConfigData, err := os.ReadFile(distConfig)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	currentConfigPath := filepath.Join(currentDir, "profile.json")
	if err := os.WriteFile(currentConfigPath, defaultConfigData, 0644); err != nil {
		return fmt.Errorf("failed to write current config: %v", err)
	}

	return nil
}

func compareFiles(file1, file2 string) (bool, error) {
	data1, err := os.ReadFile(file1)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %v", file1, err)
	}

	data2, err := os.ReadFile(file2)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %v", file2, err)
	}

	return string(data1) == string(data2), nil
}

func syncProfile(distConfig string, mergeBehavior MergeBehavior, exclude ExcludePatterns, dryRun bool) error {
	stateDirPath, err := ensureStateDir()
	if err != nil {
		return err
	}

	_, err = os.Stat(distConfig)
	if err != nil {
		return err
	}

	currentDir := filepath.Join(stateDirPath, "current")
	previousDir := filepath.Join(stateDirPath, "previous")

	// Abnormal case: reset state directory if it's invalid
	if _, err := os.Stat(currentDir); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Stat(previousDir); err == nil {
			fmt.Println("Invalid state: resetting data")
			if err := os.RemoveAll(stateDirPath); err != nil {
				return fmt.Errorf("failed to reset state directory: %v", err)
			}
			if _, err := ensureStateDir(); err != nil {
				return err
			}
		}
	}

	// TODO: applyProfile should return a new profile that only includes properties that
	// were actually changed based on the user's merge options / current settings. We can
	// still copy the distConfig but we should call it something like profile.json.orig to
	// indicate that this is just for reference
	//
	// First run: initialize current directory
	if _, err := os.Stat(currentDir); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Empty state")
		if err := applyProfile(distConfig, mergeBehavior, exclude, dryRun); err != nil {
			return err
		}
		if err := os.MkdirAll(currentDir, 0755); err != nil {
			return fmt.Errorf("failed to create current directory: %v", err)
		}
		if err := copyDistConfig(distConfig, currentDir); err != nil {
			return err
		}
		return nil
	}

	// Steady run: move current to previous and apply new config
	fmt.Println("Steady state")
	if err := os.RemoveAll(previousDir); err != nil {
		return fmt.Errorf("failed to remove previous directory: %v", err)
	}
	if err := os.Rename(currentDir, previousDir); err != nil {
		return fmt.Errorf("failed to move current to previous: %v", err)
	}
	if err := os.MkdirAll(currentDir, 0755); err != nil {
		return fmt.Errorf("failed to create current directory: %v", err)
	}
	if err := copyDistConfig(distConfig, currentDir); err != nil {
		return err
	}

	// Check if configurations differ
	currentConfig := filepath.Join(currentDir, "profile.json")
	previousConfig := filepath.Join(previousDir, "profile.json")
	identical, err := compareFiles(currentConfig, previousConfig)
	if err != nil {
		return err
	}

	if !identical {
		fmt.Println("Configurations differ -- reverting old and applying new")
		if err := revertProfile(previousConfig, exclude, dryRun); err != nil {
			return err
		}
		if err := applyProfile(currentConfig, mergeBehavior, exclude, dryRun); err != nil {
			return err
		}
	} else {
		fmt.Println("Configurations identical -- no changes required")
	}

	return nil
}
