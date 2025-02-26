package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// xfconf-query does not have a simple way to get the default value for a given property. The best we can do here
// is to start a throwaway xfconfd whose config is set to the default config provided by the distribution in
// /usr/etc/xdg (or /etc/xdg for non-atomic distros).
//
// We'll start an instance of xfconfd in a throwaway D-Bus session isolated from the user's session and call
// xfconf-query inside of that session.
func gatherDefaultPropertyValuesFromConfig(queries map[string][]string, configDir string) (map[string]map[string]string, error) {
	xfconfdPath := "/usr/lib64/xfce4/xfconf/xfconfd"
	_, err := os.Stat(xfconfdPath)
	if err != nil {
		return nil, fmt.Errorf("no xfconfd found at %s", xfconfdPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a single dbus-run-session and run xfconfd inside a shell session
	sessionCmd := fmt.Sprintf(
		"XDG_CONFIG_HOME='%s' %s & while ! xfconf-query --channel xfce4-panel --list >/dev/null 2>&1; do sleep 0.05; done; exec sh", configDir, xfconfdPath)

	logger.Debug("Launching xfconfd in its own dbus session", "cmd", sessionCmd)

	cmd := exec.CommandContext(ctx, "dbus-run-session", "--", "sh", "-c", sessionCmd)

	// Create pipes for communicating with the shell
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Start the session
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start dbus-run-session: %v", err)
	}

	// Create scanner to read stdout line-by-line
	scanner := bufio.NewScanner(stdout)

	results := make(map[string]map[string]string)

	// Query properties from each channel separately
	for channel, properties := range queries {
		results[channel] = make(map[string]string)

		for _, property := range properties {
			queryCmd := fmt.Sprintf("xfconf-query --channel %q --property %q 2>&1\n", channel, property)

			// Send query to running shell
			_, err := stdin.Write([]byte(queryCmd))
			if err != nil {
				return nil, fmt.Errorf("failed to write to dbus-run-session: %v", err)
			}

			scanner.Scan()
			queryResult := scanner.Text()
			if strings.Contains(queryResult, "does not exist on channel") {
				results[channel][property] = "" // Handle missing properties
			} else {
				results[channel][property] = queryResult
			}
		}
	}

	// Close the shell session cleanly
	stdin.Write([]byte("exit\n"))
	cmd.Wait()

	// Check if it timed out
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("timeout: xfconf-query took too long to respond")
	}

	return results, nil
}

func gatherDefaultPropertyValues(queries map[string][]string) (map[string]map[string]string, error) {
	// Special case to use the test's default values if running end-to-end-test
	if _, exists := os.LookupEnv("XFCONF_PANEL_END_TO_END_TEST"); exists {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Cannot get working directory")
		}
		return gatherDefaultPropertyValuesFromConfig(queries, cwd)
	}

	configDir := "/usr/etc/xdg"

	_, err := os.Stat(configDir)
	if err != nil {
		configDir = "/etc/xdg"
	}

	_, err = os.Stat(configDir)
	if err != nil {
		return nil, fmt.Errorf("no xdg directories available - is Xfce installed?")
	}

	return gatherDefaultPropertyValuesFromConfig(queries, configDir)
}
