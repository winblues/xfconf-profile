package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "xfconf-profile",
		Short: "A CLI tool for managing Xfconf profiles",
	}

	var applyCmd = &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply an Xfconf profile from a JSON file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Applying Xfconf profile from %s\n", args[0])
		},
	}

	var revertCmd = &cobra.Command{
		Use:   "revert [path]",
		Short: "Revert an Xfconf profile from a JSON file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Reverting Xfconf profile from %s\n", args[0])
		},
	}

	var recordCmd = &cobra.Command{
		Use:   "record",
		Short: "Record the current Xfconf profile to a JSON file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Recording current Xfconf profile")
		},
	}

	rootCmd.AddCommand(applyCmd, revertCmd, recordCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

