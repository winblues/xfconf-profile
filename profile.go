package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func applyProfile(profilePath string) error {
  fmt.Printf("Applying profile %s\n", profilePath)

	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var profile map[string]map[string]interface{}
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for channel, properties := range profile {
		for property, value := range properties {
			cmd := exec.Command("xfconf-query", "-c", channel, "--property", property, "--set", fmt.Sprintf("%v", value))

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
			}

			fmt.Printf("Executed: %s\nOutput: %s\n", cmd.String(), string(output))
		}
	}

	return nil
}

func revertProfile(profilePath string) error {
	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var profile map[string]map[string]interface{}
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for channel, properties := range profile {
		for property := range properties {
			cmd := exec.Command("xfconf-query", "-c", channel, "--reset", "--property", property)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
			}

			fmt.Printf("Executed: %s\nOutput: %s\n", cmd.String(), string(output))
		}
	}

	return nil
}

func recordProfile() {
    fmt.Printf("TODO\n")
}
