package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information overriden using ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Select a merge behavior either from the user's config or command-line flag
func chooseMergeBehavior(cfg *Config, flag string) MergeBehavior {
	if flag == "" {
		return cfg.Merge
	} else {
		parsed, err := ParseMergeBehavior(flag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return parsed
	}
}

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

			mergeFlag, _ := cmd.Flags().GetString("merge")
			mergeBehavior := chooseMergeBehavior(cfg, mergeFlag)

			distProfile, _ := cmd.Flags().GetString("profile")
			if err := syncProfile(distProfile, mergeBehavior, cfg.Exclude); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringP("profile", "p", "/usr/share/xfconf-profile/default.json", "Path to the distribution's recommended profile")
	cmd.Flags().StringP("merge", "m", "", "Set merge behavior (soft, hard, force)")
	cmd.Flags().Bool("auto", false, "Flag indicating running as a user-level systemd unit by the distribution")
	return cmd
}

func applyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply changes from a profile.json",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			mergeFlag, _ := cmd.Flags().GetString("merge")
			mergeBehavior := chooseMergeBehavior(cfg, mergeFlag)

			err := applyProfile(args[0], mergeBehavior, cfg.Exclude)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		},
	}
	cmd.Flags().StringP("merge", "m", "soft", "Set merge behavior (soft, hard, force)")
	return cmd
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
	Short: "Record changes to xfconf properties and dump them as a profile on SIGINT",
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

	rootCmd.AddCommand(applyCmd(config), revertCmd, recordCmd, syncCmd(config), versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
