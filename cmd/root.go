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
	noColor       = false
)

var rootCmd = &cobra.Command{
	Use:   "bun",
	Short: "DC/OS diagnostics bundle analysis tool",
	Long: "Bun extracts useful facts from hundreds of files in the DC/OS diagnostics bundle\n" +
		"and searches for some common problems of the DC/OS cluster.\n" +
		"\nSpecify a sub-command to run a specific check, e.g. `bun health`\n" +
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
		"print detailed output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "do not show ASCII colors")
	checks.RegisterSearchChecks()
	// Adding registered checks as commands.
	for _, c := range checks.Checks() {
		run := func(cmd *cobra.Command, args []string) {
			check := checks.GetCheck(cmd.Use)
			results := check.Run(*currentBundle)
			printReport(check, results, true)
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

func preRun(_ *cobra.Command, _ []string) {
	if currentBundle != nil {
		return
	}
	b, err := bundle.New(bundlePath)
	if err != nil {
		fmt.Printf("Cannot open a bundle: %v\n", err.Error())
		os.Exit(1)
	}
	currentBundle = &b
}

func runCheck(_ *cobra.Command, _ []string) {
	c := checks.Checks()
	sort.Slice(c, func(i, j int) bool {
		return c[i].Name < c[j].Name
	})
	ok := true
	allResults := make([]checks.Results, 0, len(c))
	for _, check := range c {
		results := check.Run(*currentBundle)
		allResults = append(allResults, results)
		printReport(check, results, false)
		if results.Status() != checks.SOK {
			ok = false
		}
	}
	printSummary(allResults)
	if !ok {
		os.Exit(1)
	}
}

// Execute starts Bun.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
