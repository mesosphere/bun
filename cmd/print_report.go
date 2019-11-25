package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/gosuri/uitable"
	"github.com/mesosphere/bun/v2/checks"
)

func printReport(c checks.Check, status checks.Status, summary string, details checks.Details) {
	au := aurora.NewAurora(!outputRedirectedToFile())
	var coloredStatus string
	switch status {
	case checks.SOK:
		coloredStatus = au.Bold(au.Green(status)).String()
	case checks.SProblem:
		coloredStatus = au.Bold(au.Red(status)).String()
	default:
		coloredStatus = au.Bold(au.Yellow(status)).String()
	}
	if status == checks.SOK && !verbose {
		return
	}
	table := uitable.New()
	table.MaxColWidth = 180
	table.Wrap = true // wrap columns
	table.AddRow("Status", coloredStatus)
	table.AddRow("Check", c.Name)
	table.AddRow("Description:", c.Description)
	table.AddRow("Summary", summary)
	if status == checks.SProblem {
		table.AddRow("Cure:", c.Cure)
	}
	table.AddRow("")
	fmt.Println(table)
}
