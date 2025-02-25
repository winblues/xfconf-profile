package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)


var applyCmd = &cobra.Command{
    Use:   "apply [path]",
    Short: "Apply changes from a profile.json",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        err := applyProfile(args[0])
        if err != nil {
          fmt.Printf("Error: %s\n", err)
        }
    },
}

var revertCmd = &cobra.Command{
    Use:   "revert [path]",
    Short: "Revert changes from a profile.json",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        err := revertProfile(args[0])
        if err != nil {
          fmt.Printf("Error: %s\n", err)
        }
    },
}

var recordCmd = &cobra.Command{
    Use:   "record",
    Short: "Record changes to xsettings and dump them as a profile on SIGINT",
    Run: func(cmd *cobra.Command, args []string) {
        recordProfile()
    },
}

var syncCmd = &cobra.Command{
    Use:   "sync",
    Short: "Sync user profile with distribution's recommended profile",
    Run: func(cmd *cobra.Command, args []string) {
        distProfile, _ := cmd.Flags().GetString("profile")
        if err := syncProfile(distProfile); err != nil {
          fmt.Fprintf(os.Stderr, "Error: %v\n", err)
          os.Exit(1)
        }
    },
}

func init() {
	syncCmd.Flags().StringP("profile", "p", "/usr/share/xfconf-profile/default.json", "Path to the distribution's recommended profile")
}

func main() {
    var rootCmd = &cobra.Command{
        Use:   "xfconf-profile",
        Short: "Tool for applying, reverting and managing Xfce profiles",
    }
    rootCmd.AddCommand(applyCmd, revertCmd, recordCmd, syncCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
