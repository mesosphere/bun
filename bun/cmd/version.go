package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version string
	Commit string
	Date string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Bun version, commit, and release date",
	Run:   printVersion,
}

func printVersion(*cobra.Command, []string) {
	fmt.Printf("Bun version: %v, commit: %v, release date: %v\n", Version, Commit, Date)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
