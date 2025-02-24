package main

import "fmt"

func applyProfile(profilePath string) {
    fmt.Printf("Applying profile %s\n", profilePath)
}

func revertProfile(profilePath string) {
    fmt.Printf("Reverting profile %s\n", profilePath)
}

func recordProfile() {
    fmt.Printf("Recording profile...\n")
}
