package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version imformation overriden using ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func syncCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync user profile with distribution's recommended profile",
		Run: func(cmd *cobra.Command, args []string) {
			auto, _ := cmd.Flags().GetBool("auto")

			if auto && !cfg.Sync.Auto {
				fmt.Println("Auto sync disabled in user config")
				return
			}

			distProfile, _ := cmd.Flags().GetString("profile")
			if err := syncProfile(distProfile); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringP("profile", "p", "/usr/share/xfconf-profile/default.json", "Path to the distribution's recommended profile")
	cmd.Flags().Bool("auto", false, "Set when running as a user-level systemd unit by the distribution")
	return cmd
}

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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("xfconf-profile version %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build date: %s\n", date)
	},
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "xfconf-profile",
		Short: "Tool for applying, reverting and managing Xfce profiles",
	}

	rootCmd.AddCommand(applyCmd, revertCmd, recordCmd, syncCmd(config), versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
