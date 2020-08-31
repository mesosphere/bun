package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mesosphere/bun/v2/tools/tasks"

	"github.com/mesosphere/bun/v2/tools/files"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Bun development tool",
	Long:  "Contains subcommands which help to add new file types and checks.",
}

func findFiles(cmd *cobra.Command, _ []string) {
	fileTypes, err := files.FindFiles(bundlePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	y, err := yaml.Marshal(&fileTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
	escape, err := cmd.Flags().GetBool("escape")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if escape {
		y = bytes.ReplaceAll(y, []byte("`"), []byte("`+ \"`\" +`"))
	}
	fmt.Println(string(y))
}

func tasksToCSV(*cobra.Command, []string) {
	err := tasks.ToCSV(currentBundle, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func init() {
	rootCmd.AddCommand(toolCmd)
	var findFileCmd = &cobra.Command{
		Use:   "find-files",
		Short: "Finds all file types in a given bundle",
		Long: "Finds all file types in a given bundle, suggests names, and" +
			" renders it in a YAML format to the stdout.",
		Run:    findFiles,
		PreRun: preRun,
	}
	findFileCmd.Flags().BoolP("escape", "e", false, "Escape back ticks for using in the files_yaml.go")
	toolCmd.AddCommand(findFileCmd)

	var tasksToCSV = &cobra.Command{
		Use:    "tasks-to-csv",
		Short:  "Represents Mesos tasks in CSV format",
		Long:   "Prints CSV table with information about Mesos tasks to stdout.",
		Run:    tasksToCSV,
		PreRun: preRun,
	}
	toolCmd.AddCommand(tasksToCSV)
}
