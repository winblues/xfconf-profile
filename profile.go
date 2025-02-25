package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os/exec"
    "github.com/fatih/color"
)

func applyProfile(profilePath string) error {
    data, err := ioutil.ReadFile(profilePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }

    var profile map[string]map[string]interface{}
    err = json.Unmarshal(data, &profile)
    if err != nil {
        return fmt.Errorf("failed to parse JSON: %v", err)
    }

    blue := color.New(color.FgHiBlue).SprintFunc()

    for channel, properties := range profile {
        for property, value := range properties {
            fmt.Printf("%s Setting %s::%s ➔ %s\n", blue("•"), channel, property, value)
            cmd := exec.Command("xfconf-query", "-c", channel, "--property", property, "--set", fmt.Sprintf("%v", value))

            output, err := cmd.CombinedOutput()
            if err != nil {
                return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
            }
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

    yellow := color.New(color.FgHiYellow).SprintFunc() 

    for channel, properties := range profile {
        for property := range properties {
            fmt.Printf("%s Resetting %s::%s\n", yellow("•"), channel, property)
            cmd := exec.Command("xfconf-query", "-c", channel, "--reset", "--property", property)

            output, err := cmd.CombinedOutput()
            if err != nil {
                return fmt.Errorf("failed to run command: %v\nOutput: %s", err, string(output))
            }
        }
    }

    return nil
}
