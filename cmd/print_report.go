package cmd

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun/v2/checks"
)

func printReport(c checks.Check) {
	printEmptyLine := false
	fmt.Printf("[%v] \"%v\" - %v\n", c.Status, c.Name, c.Summary)
	if len(c.Problems) > 0 {
		printProblems(c)
		printEmptyLine = true
	}
	if len(c.Errors) > 0 {
		printErrors(c)
		printEmptyLine = true
	}
	if verbose {
		if len(c.OKs) > 0 {
			fmt.Println("-------")
			fmt.Println("Details")
			fmt.Println("-------")
			fmt.Println(strings.Join(c.OKs, "\n"))
			printEmptyLine = true
		}
	}
	if printEmptyLine {
		fmt.Print("\n")
	}
}

func printErrors(c checks.Check) {
	fmt.Println("------")
	fmt.Println("Errors")
	fmt.Println("------")
	fmt.Println(strings.Join(c.Errors, "\n"))
}

func printProblems(c checks.Check) {
	fmt.Println("---------------")
	fmt.Println("Problem details")
	fmt.Println("---------------")
	fmt.Println(strings.Join(c.Problems, "\n"))
	if c.Cure != "" {
		printCure(c)
	}
}

func printCure(c checks.Check) {
	fmt.Println("----------------------------")
	fmt.Println("Explanation and possible fix")
	fmt.Println("----------------------------")
	fmt.Println(c.Cure)
}
