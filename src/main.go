package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Version information overriden using ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var logger *slog.Logger

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

func createSyncCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync user profile with distribution's recommended profile",
		Run: func(cmd *cobra.Command, args []string) {
			auto, _ := cmd.Flags().GetBool("auto")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			if auto && !cfg.Sync.Auto {
				fmt.Println("Auto sync disabled in user config")
				return
			}

			mergeFlag, _ := cmd.Flags().GetString("merge")
			mergeBehavior := chooseMergeBehavior(cfg, mergeFlag)

			distProfile, _ := cmd.Flags().GetString("profile")
			if err := syncProfile(distProfile, mergeBehavior, cfg.Exclude, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringP("profile", "p", "/usr/share/xfconf-profile/default.json", "Path to the distribution's recommended profile")
	cmd.Flags().StringP("merge", "m", "", "Set merge behavior (soft, hard, force)")
	cmd.Flags().Bool("dry-run", false, "Only print what would be changed")
	cmd.Flags().Bool("auto", false, "Flag indicating running as a user-level systemd unit by the distribution")

	return cmd
}

func createApplyCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply changes from a profile.json",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			mergeFlag, _ := cmd.Flags().GetString("merge")
			mergeBehavior := chooseMergeBehavior(cfg, mergeFlag)

			err := applyProfile(args[0], mergeBehavior, cfg.Exclude, dryRun)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("merge", "m", "soft", "Set merge behavior (soft, hard, force)")
	cmd.Flags().Bool("dry-run", false, "Only print what would be changed")
	return cmd
}

func createRevertCmd(cfg *Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "revert [path]",
		Short: "Revert changes from a profile.json",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			err := revertProfile(args[0], cfg.Exclude, dryRun)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().Bool("dry-run", false, "Only print what would be changed")
	return cmd
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record changes to xfconf properties and dump them as a profile",
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

func createGetDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-default <channel> <property>",
		Short: "Query xfconfd for the default value of a given property",
		Long: `Query xfconfd for the default value of a given property

      If the property does not exist, the command will return nothing on stdout and
      exit with code 127.

      Default values are found using xfconf-query on the distribution's provided Xfce settings,
      either found in /usr/etc/xdg or /etc/xdg.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			channel := args[0]
			property := args[1]

			query := map[string][]string{
				channel: {property},
			}

			defaultValues, err := gatherDefaultPropertyValues(query)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			defaultValue := defaultValues[channel][property]
			if defaultValue != "" {
				fmt.Println(defaultValue)
			} else {
				os.Exit(127)
			}
		},
	}

	return cmd
}

func initLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level})

	logger = slog.New(handler)
}

func main() {
	initLogger()

	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "xfconf-profile",
		Short: "Tool for applying, reverting and managing Xfce profiles",
	}

	applyCmd := createApplyCmd(config)
	revertCmd := createRevertCmd(config)
	syncCmd := createSyncCmd(config)
	getDefaultCmd := createGetDefaultCmd()

	rootCmd.AddGroup(&cobra.Group{ID: "profile", Title: "Profile Management"})
	applyCmd.GroupID = "profile"
	revertCmd.GroupID = "profile"
	syncCmd.GroupID = "profile"
	recordCmd.GroupID = "profile"
	rootCmd.AddCommand(applyCmd, revertCmd, syncCmd, getDefaultCmd, versionCmd, recordCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
