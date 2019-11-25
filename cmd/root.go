package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"

	"github.com/spf13/cobra"
)

var (
	bundlePath    string
	currentBundle *bundle.Bundle
	verbose       = false
)

var rootCmd = &cobra.Command{
	Use:   "bun",
	Short: "DC/OS diagnostics bundle analysis tool",
	Long: "Bun extracts useful facts from hundreds of files in the DC/OS diagnostics bundle\n" +
		"and searches for some common problems of the DC/OS cluster.\n" +
		"\nSpecify a subcommand to run a specific check, e.g. `bun health`\n" +
		"or run all the available checks by not specifying any, i.e. `bun`.\n" +
		"\nMore information is available at https://github.com/mesosphere/bun",
	PreRun: preRun,
	Run:    runCheck,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run specific check",
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while detecting a working directory: %v\n", err.Error())
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&bundlePath, "path", "p", wd,
		"path to the bundle directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"print details")
	checks.RegisterSearchChecks()
	// Adding registered checks as commands.
	for _, c := range checks.Checks() {
		run := func(cmd *cobra.Command, args []string) {
			check := checks.GetCheck(cmd.Use)
			status, summary, details := check.Run(*currentBundle)
			printReport(check, status, summary, details)
			return
		}
		var cmd = &cobra.Command{
			Use:    c.Name,
			Short:  c.Description,
			Long:   c.Description,
			PreRun: preRun,
			Run:    run,
		}
		checkCmd.AddCommand(cmd)
		checkCmd.ValidArgs = append(rootCmd.ValidArgs, cmd.Use)
		checkCmd.PreRun = preRun
	}
	rootCmd.AddCommand(checkCmd)
}

func preRun(cmd *cobra.Command, args []string) {
	if currentBundle != nil {
		return
	}
	b, err := bundle.NewBundle(bundlePath)
	if err != nil {
		fmt.Printf("Cannot find a bundle: %v\n", err.Error())
		os.Exit(1)
	}
	currentBundle = &b
}

func runCheck(cmd *cobra.Command, args []string) {
	if err := cobra.OnlyValidArgs(cmd, args); err != nil {
		fmt.Println(err.Error())
		fmt.Printf("Run '%v --help' for usage.\n", cmd.CommandPath())
		os.Exit(1)
	}
	checks := checks.Checks()
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Name < checks[j].Name
	})
	for _, check := range checks {
		status, summary, details := check.Run(*currentBundle)
		printReport(check, status, summary, details)
	}
}

// Execute starts Bun.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
