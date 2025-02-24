package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "xfconf-profile",
		Short: "A CLI tool for managing Xfce profiles",
	}

	var applyCmd = &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply changes from a profile.json",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			applyProfile(args[0])
		},
	}

	var revertCmd = &cobra.Command{
		Use:   "revert [path]",
		Short: "Revert changes from a profile.json",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			revertProfile(args[0])
		},
	}

	var recordCmd = &cobra.Command{
		Use:   "record",
		Short: "Record changes to xsettings and dump them as a profile on SIGINT",
		Run: func(cmd *cobra.Command, args []string) {
			recordProfile()
		},
	}

	rootCmd.AddCommand(applyCmd, revertCmd, recordCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
