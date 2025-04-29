/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.0.2"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the cli and exits",
	Long:  "Prints the version of the cli and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
